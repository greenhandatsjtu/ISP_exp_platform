<template>
  <v-dialog
      v-model="dialog"
      width="500"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn text dark class="ma-auto" color="indigo" v-bind="attrs" v-on="on">
        <v-icon>
          mdi-file-upload
        </v-icon>
        上传YAML
      </v-btn>
    </template>

    <v-card>
      <v-card-title class="headline grey lighten-1 font-weight-bold">
        上传YAML
      </v-card-title>
      <v-form v-model="valid">
        <v-file-input chips class="mx-2" show-size clearable label="deployment" v-model="deploy"
        />
        <v-file-input chips class="mx-2" show-size clearable label="service" v-model="svc"
        />
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
            :disabled="!(deploy&&svc)"
            @click="upload"
        >
          upload
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "UploadYamlDialog",
  inject: ['message', 'reload'],
  data: () => ({
    valid: true,
    file: null,
    dialog: false,
    deploy: null,
    svc: null,
  }),
  methods: {
    upload() {
      let data = new FormData()
      data.append('deploy', this.deploy)
      data.append('service', this.svc)
      this.$axios.post(`admin/experiment/${this.$route.params.id}/yaml`, data)
          .then(() => {
            this.message('success', '上传成功')
            this.reload()
          })
          .catch(err => {
            console.error(err)
            this.message('error', '上传失败，请检查yaml格式')
          })
    }
  }
}
</script>

<style scoped>

</style>
