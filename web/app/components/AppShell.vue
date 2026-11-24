<script setup lang="ts">
const route = useRoute()
const sidebarCookie = useCookie<boolean>("sidebar_state", { default: () => true })
const sidebarOpen = useState<boolean>("sidebar-open", () => sidebarCookie.value)

const pageTitle = usePageTitleState()
const title = computed(
  () => (pageTitle.value?.path === route.path ? pageTitle.value.title : undefined) ?? String(route.meta.title ?? "Dashboard"),
)

function updateSidebarOpen(value: boolean) {
  sidebarOpen.value = value
  sidebarCookie.value = value
}
</script>

<template>
  <SidebarProvider :open="sidebarOpen" @update:open="updateSidebarOpen">
    <AppSidebar />
    <SidebarInset>
      <header
        class="sticky top-0 z-10 flex h-14 shrink-0 items-center gap-2 border-b bg-background/85 backdrop-blur"
      >
        <div class="flex min-w-0 flex-1 items-center gap-2 px-4">
          <SidebarTrigger class="-ml-1" />
          <Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
          <Breadcrumb class="min-w-0">
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbPage class="truncate">{{ title }}</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div class="ml-auto">
            <ColorModeToggle />
          </div>
        </div>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4 sm:p-6">
        <slot />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
