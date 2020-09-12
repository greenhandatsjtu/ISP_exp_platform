<template>
  <v-snackbar
      :color="snackbar.status"
      :timeout="2000"
      app
      rounded
      top
      v-model="snackbar.open"
  >
    {{ snackbar.message }}
    <template v-slot:action="{ attrs }">
      <v-btn
          @click="snackbar.open = false"
          icon
      >
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script>
export default {
  name: "MySnackBar",
  props:['snackbar'],
  data: () => ({
  })
}
</script>

<style scoped>

</style>
