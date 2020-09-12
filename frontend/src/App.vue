<template>
  <v-app>
    <my-app-bar/>
    <v-main>
      <v-container fluid class="fill-height">
        <RouterView v-if="isRouterAlive"/>
      </v-container>
    </v-main>
    <my-snack-bar :snackbar="snackbar"/>
    <v-footer
        app
        class="justify-center"
    >
      <span>&copy; {{ new Date().getFullYear() }}</span>
    </v-footer>
  </v-app>
</template>

<script>
import HelloWorld from './components/HelloWorld';
import MyAppBar from "@/components/MyAppBar";
import MySnackBar from "@/components/MySnackBar";

export default {
  name: 'App',

  components: {
    MySnackBar,
    MyAppBar,
    HelloWorld,
  },
  provide() {
    return {
      reload: this.reload,
      message: this.message,
    }
  },
  data() {
    return {
      isRouterAlive: true,
      snackbar: {
        status: "success",
        message: null,
        open: false,

      }
    }
  },
  methods: {
    reload() {
      this.isRouterAlive = false
      this.$nextTick(function () {
        this.isRouterAlive = true
      })
    },
    message(status, message){
      this.snackbar.message = message
      this.snackbar.status = status
      this.snackbar.open = true
    }
  }
};
</script>
