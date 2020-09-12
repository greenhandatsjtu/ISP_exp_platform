<template>
  <v-card light class="mx-auto my-2">
    <v-card-title class="font-weight-bold headline grey lighten-1">实验管理区</v-card-title>
    <v-card-actions>
      <v-row class="px-2">
        <v-btn class="mx-auto my-1" v-for="(action,index) in actions" :key="index" :color="action.color">
          <v-icon>
            {{ action.icon }}
          </v-icon>
          {{ action.label }}
        </v-btn>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "ManageExperimentCard",
  props: ['experiment'],
  inject: ['message'],
  data: () => ({
    actions: [
      {label: '上传yaml', icon: 'mdi-file-upload', color: 'primary'},
      {label: '查看实验情况', icon: 'mdi-card-bulleted', color: 'primary'},
    ]
  }),
  methods: {
    enable() {
      this.$axios.get(`admin/experiment/${this.experiment.ID}/enable`)
          .then(() => {
            this.message('success', '成功启用')
            this.dialog = false
          })
    },
    disable() {
      this.$axios.get(`admin/experiment/${this.experiment.ID}/disable`)
          .then(() => {
            this.message('success', '成功禁用')
            this.dialog = false
          })
    }
  }
}
</script>

<style scoped>

</style>
