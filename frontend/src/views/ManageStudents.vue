<template>
  <v-container fluid class="fill-height">
    <v-card class="ma-auto">
      <v-card-title>
        学生表
        <v-spacer></v-spacer>
        <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search"
            single-line
            hide-details
        ></v-text-field>
        <v-spacer/>
        <add-student-dialog/>
      </v-card-title>
      <v-data-table
          :headers="headers"
          :items="students"
          :items-per-page="20"
          :search="search"
          class="elevation-10"
          item-key="student_number"
          :sort-by="['student_number']"
          show-group-by
          multi-sort
      >
        <template v-slot:item.user.online="{ item }">
          <v-icon color="green" v-if="item.user.online">mdi-check-circle</v-icon>
          <v-icon v-else>mdi-close-circle</v-icon>
        </template>
      </v-data-table>
    </v-card>
    <v-bottom-sheet v-model="sheet" inset max-width="500">
      <template v-slot:activator="{ on, attrs }">
        <v-btn fab fixed bottom right color="primary" v-bind="attrs" v-on="on">
          <v-tooltip top>
            <template v-slot:activator="{on, attrs}">
              <v-icon x-large v-on="on" v-bind="attrs">
                mdi-account-multiple-plus
              </v-icon>
            </template>
            <span>批量导入学生</span>
          </v-tooltip>
        </v-btn>
      </template>
      <v-sheet class="text-center pa-4">
        <v-btn class="warning" @click="downloadTemplate">
          <v-icon>
            mdi-download
          </v-icon>
          下载模板
        </v-btn>
        <v-form v-model="valid" ref="form">
            <v-file-input accept=".csv,.xls,.xlsx"
                          label="上传excel文件"
                          v-model="file_name">
            </v-file-input>
          <v-divider/>
          <v-row class="mx-auto">
            <v-btn
                    color="error"
                    @click="sheet=false"
            >取消
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn
                    :disabled="!valid"
                    color="info"
                    @click="submit"
            >确认
            </v-btn>
          </v-row>
        </v-form>
      </v-sheet>
    </v-bottom-sheet>
  </v-container>
</template>

<script>
  import XLSX from 'xlsx'
  import AddStudentDialog from "@/components/AddStudentDialog";
export default {
  name: "ManageStudents",
  components: {AddStudentDialog},
  inject: ['message','reload'],
  data() {
    return {
      search: '',
      students: [],
      sheet:false,
      headers: [
        {text: '学号', value: 'student_number'},
        {text: '姓名', value: 'user.name'},
        {text: '班级', value: 'class_number'},
        {text: '邮箱', value: 'user.email'},
        {text: '在线', value: 'user.online'},
      ],
      file_name: null,
      file_context: [],
      valid: false,
    }
  },
  watch: {
    file_name() {
      let file = this.file_name // 文件信息

      if (!file) {
        this.message('error', '请上传excel文件')
        return false
      } else if (!/\.(xls|xlsx)$/.test(file.name.toLowerCase())) {
        this.message('error', '请选择excel文件上传')
        return false
      }
      const fileReader = new FileReader()
      fileReader.onload = (ev) => {
        try {
          const data = ev.target.result
          const workbook = XLSX.read(data, {
            type: 'binary' // 以字符编码的方式解析
          })
          const exlName = workbook.SheetNames[0] // 取第一张表
          this.file_context = XLSX.utils.sheet_to_json(workbook.Sheets[exlName])
              .map(function (item) {
                return {
                  student_number: item['学号'],
                  class_number: item['班级'],
                  email: item['邮箱'],
                  name: item['姓名']
                }
              })
          for (let item in this.file_context) {
            this.file_context[item].student_number = this.file_context[item].student_number.toString()
          }
          console.log(this.file_context)
        } catch (e) {
          this.message('error', '异常错误')
          return false
        }
      }
      fileReader.readAsBinaryString(file)
    }
  },
  methods: {
    getStudents() {
      this.$axios.get('admin/student')
          .then(({data}) => {
            this.students = data.data
          })
          .catch(err => {
            console.log(err)
            this.message('error', '获取学生错误')
          })
    },
    submit() {
      let that = this

      this.$axios.post('admin/students', that.file_context).then(() => {
        this.message('success', '添加成功')
        this.reload()
      }).catch(err => {
        console.log(err)
        this.message('error', '格式错误，请检查格式')
      })
    },
    downloadTemplate(){
      window.open(`api/admin/users/template`)
    }
  },
  created() {
    this.getStudents()
  }
}
</script>

<style scoped>

</style>
