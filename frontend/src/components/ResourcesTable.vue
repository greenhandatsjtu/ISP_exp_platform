<template>
  <v-container fluid class="fill-height">
    <v-card class="ma-auto">
      <v-card-title class="headline font-weight-bold">
        实验资源表
        <v-spacer></v-spacer>
        <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search"
            class="mx-2"
            single-line
            clearable
            hide-details
        ></v-text-field>
      </v-card-title>
      <v-data-table
          :headers="headers"
          :items="resources"
          :items-per-page="10"
          :search="search"
          class="elevation-20"
          item-key="ID"
          multi-sort
      >
        <template v-slot:item.CreatedAt="{ item }">
          {{ new Date(item.CreatedAt).toLocaleString() }}
        </template>
        <template v-slot:item.user.online="{ item }">
          <v-icon color="green" v-if="item.user.online">mdi-check-circle</v-icon>
          <v-icon v-else>mdi-close-circle</v-icon>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{on, attrs}">
              <v-icon
                  color="error"
                  v-on="on"
                  v-bind="attrs"
                  @click="endExp(item)"

              >
                mdi-flask-empty-off-outline
              </v-icon>
            </template>
            <span>结束实验</span>
          </v-tooltip>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>

<script>
export default {
  name: "ResourcesTable",
  inject: ['message', 'reload'],
  props: ['scope'],
  data() {
    return {
      search: '',
      resources: [],
      dialog: false,
      deleteItem: {},
    }
  },
  methods: {
    getResources() {
      let url = ''
      switch (this.scope) {
        case 'all':
          url = 'admin/resources/all'
          break
        case 'course':
          url = `admin/resources/course/${this.$route.params.id}`
          break
        case 'experiment':
          url = `admin/resources/experiment/${this.$route.params.id}`
          break
        default:
          return
      }
      this.$axios.get(url)
          .then(({data}) => {
            this.resources = data.data
          })
          .catch(err => {
            console.log(err)
            this.message('error', '获取实验资源错误')
          })
    },
    editItem(item) {
      this.editedItem = item
      let currentRoles = item.roles.map(x => x.description)
      this.unselected = []
      for (let role of this.roles) {
        if (!currentRoles.includes(role)) {
          this.unselected.push(role)
        }
      }
      this.dialog = true
    },
    endExp(item) {
      this.$axios.delete(`admin/experiment/${item.experiment_id}/user/${item.user_id}`)
          .then(() => {
            this.message('info', '成功结束实验')
            this.reload()
          })
          .catch(() => {
            this.message('error', '结束实验失败！')
          })
    },
  },
  computed: {
    headers() {
      return [
        {text: '姓名', value: 'user.name'},
        {text: '实验名', value: 'experiment.name'},
        {text: '实验开始时间', value: 'CreatedAt'},
        {text: '在线', value: 'user.online'},
        {text: '端口', value: 'port'},
        {text: '行为', value: 'actions', sortable: false,}
      ]
    },
  },
  created() {
    this.getResources()
  },
}
</script>

<style scoped>

</style>
