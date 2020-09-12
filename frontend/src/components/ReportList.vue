<template>
  <v-card>
    <v-card-title class="headline font-weight-bold">
      提交列表
      <v-tooltip top>
        <template v-slot:activator="{ on, attrs }">
          <v-btn v-on="on" v-bind="attrs" fab class="mx-2 primary" :disabled="upload.length===0" @click="downloadReports">
          <v-icon>
            mdi-folder-download
          </v-icon>
          </v-btn>
        </template>
        <span>打包下载</span>
      </v-tooltip>
      <v-spacer/>
      <v-text-field
          class="mx-2 my-0"
          v-model="search"
          append-icon="mdi-magnify"
          label="搜索"
          single-line
          hide-details
      ></v-text-field>
    </v-card-title>
    <v-data-table
        v-model="upload"
        :search="search"
        class="ma-2 elevation-10"
        :headers="headers"
        :items="allStudents"
        :items-per-page="10"
        show-group-by
        item-key="student_number"
        :sort-by="['student_number']"
    >
      <template v-slot:item.upload="{ item }">
        <v-icon color="green" v-if="item.upload">mdi-check-circle</v-icon>
        <v-icon v-else>mdi-close-circle</v-icon>
      </template>
      <template v-slot:item.upload_at="{ item }">
        {{ getUploadTime(item) }}
      </template>
      <template v-slot:item.action="{ item }">
        <v-btn icon @click="download(item)" :disabled="!item.upload">
          <v-icon color="success">mdi-download</v-icon>
        </v-btn>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
export default {
  name: "ReportList",
  data() {
    return {
      allStudents: [],
      search: '',
      headers: [
        {text: '学号', value: 'student_number'},
        {text: '班级', value: 'class_number'},
        {text: '姓名', value: 'user.name'},
        {text: '提交', value: 'upload'},
        {text: '最后提交时间', value: 'upload_at'},
        {text: '下载', value: 'action'}
      ],
      upload: [],
      HasUpload: [],
    }
  },
  inject: ['message'],
  methods: {
    getStudent() {
      this.$axios.get(`/admin/experiment/${this.$route.params.id}/student`)
          .then(({data}) => {
            this.allStudents = data.data
          }).catch((err) => {
        this.message("error", "获取学生出错")
      })
    },
    getReport() {
      this.$axios.get(`/admin/experiment/${this.$route.params.id}/reports`)
          .then(({data}) => {
            this.upload = data.data
            for (let item of this.upload) {
              this.HasUpload.push(item.student.ID)
            }
            for (let item of this.allStudents) {
              item.upload = this.HasUpload.includes(item.ID);
            }
          }).catch((err) => {
        this.message("error", "获取报告出错")
      })
    },
    download(item) {
      window.open(`api/admin/experiment/${this.$route.params.id}/student/${item.ID}/report`)
    },
    downloadReports() {
      window.open(`api/admin/experiment/${this.$route.params.id}/reports/download`)
    },
    getUploadTime(item){
      for(let report of this.upload){
        if(report.student_id===item.ID){
          return new Date(report.UpdatedAt).toLocaleString()
        }
      }
    }
  },
  mounted() {
    this.getReport()
  },
  created() {
    this.getStudent()

  }
}
</script>

<style scoped>

</style>
