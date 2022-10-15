/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import "@/css/main.css"
import { createPinia } from "pinia"
import { createApp } from "vue"
import { RouterView } from "vue-router/auto"
import { client } from "@/lib/api"
import router from "@/lib/core/router"
import urql from "@urql/vue"

const app = createApp(RouterView)

app.use(urql, client)
app.use(createPinia())
app.use(router)
app.mount("#app")
