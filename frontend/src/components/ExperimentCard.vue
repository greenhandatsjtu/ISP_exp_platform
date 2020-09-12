<template>
  <v-col :cols="12" :sm="4" class="ma-auto">
    <v-hover v-slot:default="{hover}">
      <v-card
          rounded
          :elevation="hover?20:5"
          class="mb-1"
      >
        <v-img :src="require('@/assets/exp.jpg')" max-height="200" contain/>
        <v-card-title v-text="experiment.name" class="font-weight-bold"></v-card-title>
        <v-card-actions>
          <v-btn
              text
              color="info"
              :to="{name: admin?'AdminExperimentInfo':'ExperimentInfo', params: {id: experiment.ID}}"
          >
            进入实验
          </v-btn>
          <v-spacer></v-spacer>
          <v-dialog
              light
              v-model="dialog"
              max-width="300"
          >
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                  color="error"
                  icon
                  v-if="admin"
                  v-bind="attrs"
                  v-on="on"
              >
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </template>

            <v-card>
              <v-card-title class="headline grey lighten-2">
                确定删除实验？
              </v-card-title>
              <v-divider/>
              <v-card-actions>
                <v-btn color="primary" @click="dialog=false">
                  取消
                </v-btn>
                <v-spacer/>
                <v-btn color="warning" @click="confirmDelete">
                  确定
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-card-actions>
      </v-card>
    </v-hover>
  </v-col>
</template>

<script>
export default {
  name: "ExperimentCard",
  props: ['experiment', 'admin','reload'],
  inject:['reload','message'],
  data: () => ({
    dialog: false
  }),
  methods: {
    confirmDelete() {
      this.$axios.get(`admin/experiment/${this.experiment.ID}/delete`)
          .then(() => {
            this.message('info', '成功删除实验')
            this.dialog = false
            this.reload()
          })
          .catch(()=>{
            this.message('warning', '删除实验失败')
          })
    }
  }
}
</script>

<style scoped>

</style>
