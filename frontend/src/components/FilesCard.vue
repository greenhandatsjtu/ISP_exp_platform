<template>
  <v-card light>
    <v-card-title class="font-weight-bold headline grey lighten-1">相关文件</v-card-title>
    <v-divider/>
    <v-card-subtitle v-if="files.length===0">暂无</v-card-subtitle>
    <v-list dense>
      <v-list-item v-for="(file, index) in files" :key="index">
        <v-list-item-content>
          <v-list-item-title>{{ file.file_name }}</v-list-item-title>
        </v-list-item-content>
        <v-list-item-icon>
          <v-btn icon @click="download(file.file_name)">
            <v-icon color="success">mdi-download</v-icon>
          </v-btn>
        </v-list-item-icon>
        <v-list-item-icon v-if="admin">
          <v-btn icon @click="removeFile(file.file_name)">
            <v-icon color="error">mdi-delete</v-icon>
          </v-btn>
        </v-list-item-icon>
      </v-list-item>
    </v-list>
    <v-card-actions>
      <v-btn :disabled="files.length===0" text color="primary" @click="downloadAll">
        <v-icon large>mdi-folder-download</v-icon>
        下载压缩包
      </v-btn>
      <v-spacer/>
      <upload-dialog v-if="admin" class="mx-auto my-2" title="文件" :max="50" :upload-file="uploadFile"/>
    </v-card-actions>
  </v-card>
</template>

<script>
import UploadDialog from "@/components/UploadDialog";

export default {
  name: "FilesCard",
  components: {UploadDialog},
  props: ['files', 'admin'],
  inject: ['message', 'reload'],
  methods: {
    download(fileName) {
      window.open(`api/experiment/${this.$route.params.id}/doc/${fileName}`)
    },
    downloadAll() {
      window.open(`api/experiment/${this.$route.params.id}/docs`)
    },
    removeFile(fileName) {
      this.$axios.delete(`admin/experiment/${this.$route.params.id}/doc/${fileName}`)
          .then(() => {
            this.message('success', `成功删除文件：${fileName}`)
            this.reload()
          })
          .catch(err => {
            console.error(err)
            this.message('error', "删除失败")
          })
    },
    uploadFile(file) {
      let data = new FormData();
      data.append('file', file)
      this.$axios.post(`admin/experiment/${this.$route.params.id}/doc`, data)
          .then(() => {
            this.message('success', '上传成功')
            this.reload()
          })
          .catch(err => {
            console.error(err)
            this.message('error', '上传失败')
          })
    },
  }
}
</script>

<style scoped>

</style>
