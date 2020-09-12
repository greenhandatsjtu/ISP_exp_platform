<template>
  <v-container fluid class="fill-height ma-0">
    <v-row>
      <v-alert class="ma-0" v-if="isDoing" border="left" dismissible elevation="10" type="info" dense>
        正在进行实验，端口号：<strong>{{info.port}}</strong>
      </v-alert>
    </v-row>
    <v-row>
      <v-card-text><span class="font-weight-bold">实验说明：</span>{{ experiment.assignment }}</v-card-text>
    </v-row>
    <v-row>
      <v-card-text><span class="font-weight-bold">实验时间：</span>{{ new Date(experiment.time).toLocaleString() }}
      </v-card-text>
    </v-row>
    <v-btn @click="startExp" v-if="!isDoing" :disabled="!(experiment.upload&&(admin||experiment.enable))"
           class="ma-auto" color="warning">
      <v-icon>
        mdi-flask-empty-outline
      </v-icon>
      开启实验
    </v-btn>
    <v-btn @click="endExp" v-else class="ma-auto" color="error">
      <v-icon>
        mdi-close-outline
      </v-icon>
      结束实验
    </v-btn>
  </v-container>
</template>

<script>
export default {
  name: "ExperimentInfoZone",
  props: ['experiment', 'admin'],
  inject: ['message', 'reload'],
  data() {
    return {
      isDoing: true,
      info: {
        port: 0,
      },
    }
  },
  methods: {
    startExp() {
      this.$axios.get(`${this.admin ? 'admin/' : ''}experiment/${this.$route.params.id}/start`)
          .then(({data}) => {
            this.message('info', `成功开启，端口号: ${data.data.port}`)
            this.reload()
          })
          .catch(() => {
            this.message('error', '开启失败，请管理员检查yaml')
          })
    },
    endExp() {
      this.$axios.get(`experiment/${this.$route.params.id}/end`)
          .then(() => {
            this.message('info', '成功结束实验')
            this.reload()
          })
          .catch(() => {
            this.message('error', '结束实验失败！')
          })
    }
  },
  created() {
    this.$axios.get(`experiment/${this.$route.params.id}/status`)
        .then(({data}) => {
          this.info = data.data
          this.isDoing = this.info.port!==0
        })
        .catch(() => {
          this.message('error', '查询实验状态失败')
          this.isDoing = true
        })
  }
}
</script>

<style scoped>

</style>
