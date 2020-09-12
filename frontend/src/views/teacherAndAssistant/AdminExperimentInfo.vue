<template>
  <v-container fluid class="fill-height">
    <v-row>
      <v-tooltip top>
        <template v-slot:activator="{ on, attrs }">
          <v-btn icon color="warning" class="mr-6" v-bind="attrs" v-on="on" @click.stop="openDialog">
            <v-icon x-large>
              mdi-square-edit-outline
            </v-icon>
          </v-btn>
        </template>
        <span>编辑实验信息</span>
      </v-tooltip>
      <h1 class="text-truncate">{{ experiment.name }}</h1>
      <v-btn class="ma-2" color="primary" v-if="!experiment.enable" :disabled="!experiment.upload"
             @click="enable">
        <v-icon>
          mdi-check-bold
        </v-icon>
        启用实验
      </v-btn>
      <v-btn class="ma-2" color="warning" v-else @click="disable">
        <v-icon>
          mdi-close-circle
        </v-icon>
        禁用实验
      </v-btn>
    </v-row>
    <v-row>
      <v-col>
        <experiment-info-zone :experiment="experiment" :admin="true"/>
      </v-col>
      <v-col>
        <yaml-card :experiment="experiment"/>
      </v-col>
      <v-col>
        <files-card :admin="true" :files="experiment.files"/>
      </v-col>
    </v-row>
    <v-expansion-panels focusable accordion popout>
      <v-expansion-panel>
        <v-expansion-panel-header>查看报告提交情况</v-expansion-panel-header>
        <v-expansion-panel-content>
          <ReportList/>
        </v-expansion-panel-content>
      </v-expansion-panel>
      <v-expansion-panel>
        <v-expansion-panel-header>查看实验情况</v-expansion-panel-header>
        <v-expansion-panel-content>
          <resources-table scope="experiment"/>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
    <v-dialog
        v-model="dialog"
        max-width="400"
    >
      <v-card class="px-4">
        <v-card-title>
          编辑实验
        </v-card-title>
        <v-form
            ref="form"
            v-model="valid">
          <v-text-field
              v-model="edit.name"
              label="实验名"
              required
              :rules="nameRules"
              counter
          ></v-text-field>
          <label>
            实验日期
            <input required type="datetime-local" v-model="edit.time"/>
          </label>
          <v-textarea
              v-model="edit.assignment"
              label="实验说明"
              :rules="assignmentRules"
              auto-grow
              counter
              filled
          ></v-textarea>
        </v-form>
        <v-card-actions>
          <v-btn color="primary" :disabled="!valid" @click="submit">
            确定
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import FilesCard from "@/components/FilesCard";
import ManageExperimentCard from "@/components/ManageExperimentCard";
import UploadDialog from "@/components/UploadDialog";
import ExperimentInfoZone from "@/components/ExperimentInfoZone";
import UploadYamlDialog from "@/components/UploadYamlDialog";
import ResourcesTable from "@/components/ResourcesTable";
import YamlCard from "@/components/YamlCard";
import ReportList from "../../components/ReportList";

export default {
  name: "ExperimentInfo",
  inject: ['reload', 'message'],
  components: {
    YamlCard,
    ResourcesTable, UploadYamlDialog, ExperimentInfoZone, UploadDialog, ManageExperimentCard, FilesCard, ReportList
  },
  data() {
    return {
      valid: true,
      dialog: false,
      file: null,
      grade: null,
      experiment: {
        files: [],
      },
      edit: {
        name: '',
        assignment: '',
        time: null,
      },
      nameRules: [
        v => !!v || '实验名未填！',
        v => v.length <= 50 || '实验名过长！',
      ],
      assignmentRules: [
        v => v.length <= 100 || '实验说明过长！',
      ],
    }
  },
  methods: {
    fetchExperiment() {
      this.$axios.get(`experiment/${this.$route.params.id}`)
          .then(({data}) => {
            this.experiment = data.data
          })
    },
    openDialog() {
      this.dialog = true
      this.edit.name = this.experiment.name
      this.edit.assignment = this.experiment.assignment
    },
    submit() {
      let data = {
        name: this.edit.name,
        assignment: this.edit.assignment,
      }
      if (this.edit.time) {
        data['time'] = new Date(this.edit.time);
      }
      this.$axios.post(`admin/experiment/${this.$route.params.id}/update`, data)
          .then(() => {
            this.message('success', '成功修改')
            this.dialog = false
            this.reload()
          })
          .catch(err => {
            console.log(err)
            this.message('warning', '修改失败')
          })
    },
    enable() {
      this.$axios.get(`admin/experiment/${this.experiment.ID}/enable`)
          .then(() => {
            this.message('success', '成功启用')
            this.dialog = false
            this.reload()
          })
    },
    disable() {
      this.$axios.get(`admin/experiment/${this.experiment.ID}/disable`)
          .then(() => {
            this.message('success', '成功禁用')
            this.dialog = false
            this.reload()
          })
    }
  },
  created() {
    this.fetchExperiment()
  }

}
</script>

<style scoped>

</style>
