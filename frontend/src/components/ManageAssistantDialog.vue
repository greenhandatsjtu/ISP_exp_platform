<template>
  <v-dialog
      v-model="dialog"
      max-width="600"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn icon v-on="on" v-bind="attrs">
        <v-icon color="teal">
          mdi-account-cog
        </v-icon>
      </v-btn>
    </template>

    <v-card>
      <v-card-title>
        管理助教
      </v-card-title>
      <v-data-table
          v-model="selected"
          class="ma-2 elevation-10"
          :headers="headers"
          :items="allAssistants"
          :items-per-page="5"
          item-key="ID"
          show-select
      >
      </v-data-table>
      <v-card-actions>
        <v-btn color="primary" @click="updateAssistants">
          确定
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "ManageAssistantDialog",
  inject: ['message','reload'],
  data: function () {
    return {
      dialog: false,
      headers: [
        {text: '姓名', value: 'user.name'},
        {text: '邮箱', value: 'user.email'},
      ],
      allAssistants: [],
      selected: []
    }
  },
  created() {
    this.$axios.get('admin/assistants')
        .then(({data}) => {
          this.allAssistants = data.data
        })
        .catch(err => {
          console.error(err)
          this.message('warning', '获取助教失败')
        })
  },
  mounted() {
    this.$axios.get(`course/${this.$route.params.id}/assistants`)
        .then(({data}) => {
          this.selected = data.data
        })
        .catch(err => {
          console.error(err)
          this.message('warning', '获取已分配助教失败')
        })
  },
  methods: {
    updateAssistants() {
      this.$axios.post(`admin/course/${this.$route.params.id}/assistants`, {
        assistants: this.selected.map(x => x.ID)
      })
          .then(() => {
            this.message('success', '分配成功')
            this.dialog = false
            this.reload()
          })
          .catch(() => {
            this.message('warning', '分配失败')
          })
    }
  }
}
</script>

<style scoped>

</style>
