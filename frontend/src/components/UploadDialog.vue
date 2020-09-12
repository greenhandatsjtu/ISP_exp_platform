<template>
  <v-dialog
      v-model="dialog"
      width="500"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn text class="ma-auto" color="teal" v-bind="attrs" v-on="on">
        <v-icon>
          mdi-file-upload
        </v-icon>
        上传{{ title }}
      </v-btn>
    </template>

    <v-card>
      <v-card-title class="headline grey lighten-1 font-weight-bold">
        上传{{ title }}
      </v-card-title>
      <v-form v-model="valid">
        <v-file-input class="mx-2" chips show-size clearable :label="title" v-model="file"
                      :rules="rules"/>
      </v-form>
      <v-divider></v-divider>
      <v-card-actions>
        <v-btn
            text
            @click="dialog=false"
        >
          cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            text
            color="primary"
            :disabled="!valid"
            @click="uploadFile(file)"
        >
          upload
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "UploadDialog",
  props: ['uploadFile', 'title', 'max'],
  data: () => ({
    valid: true,
    file: null,
    dialog: false,
    rules: [],
  }),
  mounted() {
    this.rules = [
      v => !!v || '请选择文件',
      value => value && value.size < this.max*1000000 || `文件不能超过${this.max}MB!`,
    ]
  }
}
</script>

<style scoped>

</style>
