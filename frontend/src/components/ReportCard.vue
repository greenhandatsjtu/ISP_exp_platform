<template>
  <v-card light>
    <v-card-title class="font-weight-bold headline grey lighten-1">实验报告</v-card-title>
    <v-divider/>
    <v-card-subtitle v-if="!report">暂无</v-card-subtitle>
    <v-list dense v-else>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title>{{ report }}</v-list-item-title>
        </v-list-item-content>
        <v-list-item-icon>
          <v-btn icon @click="download()">
            <v-icon color="success">mdi-download</v-icon>
          </v-btn>
        </v-list-item-icon>
      </v-list-item>
    </v-list>
    <v-card-actions>
      <upload-dialog title="报告" :max="10" :upload-file="uploadFile"/>
    </v-card-actions>
  </v-card>
</template>

<script>
import UploadDialog from "@/components/UploadDialog";

export default {
  name: "ReportCard",
  components: {UploadDialog},
  inject: ['message', 'reload'],
  data() {
    return {
      report: null,
    }
  },
  methods: {
    download() {
      window.open(`api/student/experiment/${this.$route.params.id}/report/download`)
    },
    getReport() {
      this.$axios.get(`student/experiment/${this.$route.params.id}/report`)
          .then(({data}) => {
            this.report = data.data.file_name
          }).catch(() => {

      })
    },
    uploadFile(file) {
      let data = new FormData();
      data.append('file', file)
      this.$axios.post(`student/experiment/${this.$route.params.id}/report`, data)
          .then(() => {
            this.message('success', "上传成功")
            this.reload()
          })
          .catch(err => {
            console.error(err)
            this.message('error', "上传失败")
          })
    },
  },
  created() {
    this.getReport()
  }
}
</script>

<style scoped>

</style>
