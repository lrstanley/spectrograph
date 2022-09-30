/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import "nprogress/nprogress.css"
import { useNProgress } from "@vueuse/integrations/useNProgress"

export const { isLoading: loadingBar } = useNProgress()
