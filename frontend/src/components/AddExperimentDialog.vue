<template>
  <v-dialog
      v-model="dialog"
      max-width="400"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn class="ma-auto" color="indigo" dark v-bind="attrs" v-on="on">
        <v-icon>
          mdi-calendar-plus
        </v-icon>
        新增实验
      </v-btn>
    </template>

    <v-card class="px-4">
      <v-card-title>
        新增实验
      </v-card-title>
      <v-form
          ref="form"
          v-model="valid">
        <v-text-field
            v-model="experiment.name"
            label="实验名"
            required
            :rules="nameRules"
            counter
        ></v-text-field>
        <label>
          实验日期
          <input required type="datetime-local" v-model="experiment.time"/>
        </label>
        <v-textarea
            v-model="experiment.assignment"
            label="实验说明"
            :rules="assignmentRules"
            auto-grow
            counter
            filled
        ></v-textarea>
      </v-form>
      <v-card-actions>
        <v-btn color="primary" :disabled="!valid" @click="addExperiment">
          确定
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "AddExperimentDialog",
  inject: ['reload', 'message'],
  data: function () {
    return {
      dialog: false,
      valid: true,
      experiment: {
        name: '',
        time: '',
        assignment: '',
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
    addExperiment() {
      if(this.experiment.time){
        this.experiment.time = new Date(this.experiment.time);
      }else {
        this.experiment.time = new Date();
      }
      this.$axios.post(`admin/course/${this.$route.params.id}/experiment`, this.experiment)
          .then(() => {
            this.message('success', '添加成功')
            this.dialog = false
            this.reload()
          })
          .catch(err => {
            console.log(err)
            this.message('warning', '添加失败')
          })
    }
  }
}
</script>

<style scoped>

</style>
