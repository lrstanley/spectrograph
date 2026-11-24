import { type ChildProcess, spawn } from "node:child_process"
import { randomBytes } from "node:crypto"
import { createServer } from "node:net"
import { fileURLToPath } from "node:url"

const repoRoot = fileURLToPath(new URL("../..", import.meta.url))

export type SpectrographServer = {
  url: string
  stop: () => Promise<void>
}

async function freePort(): Promise<number> {
  const server = createServer()
  await new Promise<void>((resolve, reject) => {
    server.once("error", reject)
    server.listen(0, "127.0.0.1", () => resolve())
  })
  const addr = server.address()
  if (addr == null || typeof addr === "string") {
    server.close()
    throw new Error("no free tcp port")
  }
  const port = addr.port
  await new Promise<void>((resolve) => server.close(() => resolve()))
  return port
}

async function waitHealthy(url: string, child: ChildProcess): Promise<void> {
  const logs: string[] = []
  const onLog = (buf: Buffer) => logs.push(buf.toString())
  child.stdout?.on("data", onLog)
  child.stderr?.on("data", onLog)

  const deadline = Date.now() + 180_000
  while (Date.now() < deadline) {
    if (child.exitCode != null) {
      throw new Error(`httpserver exited ${child.exitCode}: ${logs.join("")}`)
    }
    try {
      const res = await fetch(`${url}/-/healthy`)
      if (res.ok) return
    } catch {
      // Compile or bind still in progress.
    }
    await new Promise((resolve) => setTimeout(resolve, 250))
  }
  throw new Error(`httpserver did not become healthy: ${logs.join("")}`)
}

// Starts the backend (with the embedded frontend, so build it first with
// "task node-build"), with the Discord worker disabled and dummy OAuth
// credentials. Only flows that don't need a Discord login can be tested.
export async function startSpectrograph(databaseURL: string): Promise<SpectrographServer> {
  const port = await freePort()
  const url = `http://localhost:${port}`

  const child = spawn("go", ["run", "./cmd/httpserver", "--no-dot-env"], {
    cwd: repoRoot,
    env: {
      PATH: process.env.PATH,
      HOME: process.env.HOME,
      GOPATH: process.env.GOPATH ?? "",
      GOCACHE: process.env.GOCACHE ?? "",
      GOMODCACHE: process.env.GOMODCACHE ?? "",
      DATABASE_URL: databaseURL,
      HTTP_BIND_ADDR: `127.0.0.1:${port}`,
      HTTP_BASE_URL: url,
      HTTP_VALIDATION_KEY: randomBytes(64).toString("hex"),
      HTTP_ENCRYPTION_KEY: randomBytes(32).toString("hex"),
      DISCORD_CLIENT_ID: "123456",
      DISCORD_CLIENT_SECRET: "e2e",
      WORKER_ENABLED: "false",
    },
    stdio: ["ignore", "pipe", "pipe"],
    detached: true,
  })

  const stop = async () => {
    if (child.pid != null) {
      try {
        process.kill(-child.pid, "SIGTERM")
      } catch {
        child.kill("SIGTERM")
      }
    }
    await new Promise<void>((resolve) => {
      const timer = setTimeout(() => {
        if (child.pid != null) {
          try {
            process.kill(-child.pid, "SIGKILL")
          } catch {
            child.kill("SIGKILL")
          }
        }
        resolve()
      }, 5000)
      child.once("exit", () => {
        clearTimeout(timer)
        resolve()
      })
    })
  }

  try {
    await waitHealthy(url, child)
  } catch (err) {
    await stop()
    throw err
  }
  return { url, stop }
}
