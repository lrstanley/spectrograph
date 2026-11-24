<script setup lang="ts">
import { toast } from "vue-sonner"

const open = defineModel<boolean>("open", { required: true })

const { $api } = useNuxtApp()
const { clear } = useSelf()
const pending = ref(false)
const error = ref("")

async function confirm() {
  pending.value = true
  error.value = ""
  try {
    await $api("/self", { method: "DELETE" })
    clear()
    clearNuxtData("my-guilds")
    open.value = false
    toast.success("Account deleted")
    await navigateTo("/")
  } catch (err) {
    error.value = apiErrorMessage(err)
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <AlertDialog v-model:open="open">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>Delete your account?</AlertDialogTitle>
        <AlertDialogDescription>
          This deletes your Spectrograph account and signs you out. It doesn't remove the bot from your servers, since
          other admins may still use it there: pause it on each server first, or remove the bot in Discord.
        </AlertDialogDescription>
      </AlertDialogHeader>
      <Alert v-if="error" variant="destructive">
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>
      <AlertDialogFooter>
        <AlertDialogCancel :disabled="pending">Keep account</AlertDialogCancel>
        <Button variant="destructive" :disabled="pending" @click="confirm">
          <Spinner v-if="pending" data-icon="inline-start" />
          Delete account
        </Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
