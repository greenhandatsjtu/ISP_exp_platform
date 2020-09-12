<template>
  <v-container fluid class="fill-height">
    <v-row>
      <h1 class="text-truncate">{{ experiment.name }}</h1>
    </v-row>
    <v-row>
      <v-col>
        <v-row v-if="!!grade">
          <v-card-title class="headline"><span class="font-weight-bold">成绩：</span>{{ grade }}</v-card-title>
        </v-row>
        <experiment-info-zone :experiment="experiment"/>
      </v-col>
      <v-col>
        <report-card/>
      </v-col>
      <v-col>
        <files-card :files="experiment.files"/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import FilesCard from "@/components/FilesCard";
import Roles from "@/components/Roles";
import ManageExperimentCard from "@/components/ManageExperimentCard";
import UploadDialog from "@/components/UploadDialog";
import ExperimentInfoZone from "@/components/ExperimentInfoZone";
import ReportCard from "../../components/ReportCard";

export default {
  name: "ExperimentInfo",
  components: {ReportCard,ExperimentInfoZone, UploadDialog, ManageExperimentCard, FilesCard},
  inject:['message','reload'],
   data() {
    return {
      grade: null,
      experiment: {
        files: [],
      },
    }
  },
  methods: {
    fetchExperiment() {
      this.$axios.get(`experiment/${this.$route.params.id}`)
          .then(({data}) => {
            this.experiment = data.data
          })
    },
    getGrade() {
      this.$axios.get(`student/experiment/${this.$route.params.id}/grade`)
          .then(({data}) => {
            this.grade = data.data.grade
          })
          .catch(err => {
          })
    },
  },
  created() {
    this.fetchExperiment()
    // this.getGrade()
  }

}
</script>

<style scoped>

</style>
