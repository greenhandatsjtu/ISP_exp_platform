<template>
  <v-container fluid class="fill-height">
    <v-card class="ma-auto" light>
      <v-card-title class="headline font-weight-bold blue lighten-3">
        用户表
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
        <v-select single-line prepend-icon="mdi-account-outline" class="mx-1" hide-details :items="chRoles" v-model="role"
                      label="搜索角色" clearable></v-select>
      </v-card-title>
      <v-data-table
          :headers="headers"
          :items="users"
          :items-per-page="10"
          :search="search"
          class="elevation-20"
          item-key="ID"
          multi-sort
      >
        <template v-slot:item.online="{ item }">
          <v-icon color="green" v-if="item.online">mdi-check-circle</v-icon>
          <v-icon v-else>mdi-close-circle</v-icon>
        </template>
        <template v-slot:item.roles="{ item }">
          <v-chip class="mr-1 mb-1 font-weight-bold" dark small v-for="(role,index) in item.roles" :key="index"
                  :color="rolesMap[role.description][1]">{{
              rolesMap[role.description][0]
            }}
          </v-chip>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{on, attrs}">
              <v-icon
                  class="mr-1"
                  color="primary"
                  v-on="on"
                  v-bind="attrs"
                  @click="editItem(item)"
              >
                mdi-pencil
              </v-icon>
            </template>
            <span>编辑角色和密码</span>
          </v-tooltip>
          <v-tooltip top>
            <template v-slot:activator="{on, attrs}">
              <v-icon
                  color="error"
                  v-on="on"
                  v-bind="attrs"
                  @click="openDeleteDialog(item)"

              >
                mdi-delete
              </v-icon>
            </template>
            <span>删除用户</span>
          </v-tooltip>
        </template>
      </v-data-table>
    </v-card>
    <v-dialog
        light
        max-width="500"
        v-model="dialog"
    >
      <v-card>
        <v-card-title class="grey lighten-2">
          编辑{{ editedItem.name }}的角色和密码
        </v-card-title>
        <v-row class="mx-2 my-1">
          <v-chip @click:close="openRevokeDialog(role.description)" close close-icon="mdi-minus"
                  class="ma-2 font-weight-bold" dark small
                  v-for="(role,index) in editedItem.roles" :key="index"
                  :color="rolesMap[role.description][1]">{{
              rolesMap[role.description][0]
            }}
          </v-chip>
        </v-row>
        <v-row class="mx-2 my-1">
          <v-chip @click:close="openAddDialog(role)" close close-icon="mdi-plus" class="ma-2 font-weight-bold" dark small
                  v-for="(role,index) in unselected" :key="index"
                  :color="rolesMap[role][1]">{{
              rolesMap[role][0]
            }}
          </v-chip>
        </v-row>
        <v-card-text>
          注：学生角色不能分配、撤回（请到学生管理页面），助教角色只能分配给学生
        </v-card-text>
        <v-form v-model="valid" ref="form">
          <v-row class="mx-2 my-1">
            <v-text-field :rules="passwordRules" single-line v-model="password"/>
            <v-btn color="warning" @click="changePassword" :disabled="!valid">
              修改密码
            </v-btn>
          </v-row>
        </v-form>
        <v-card-actions>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
        light
        v-model="deleteDialog"
        max-width="400"
    >
      <v-card>
        <v-card-title class="grey lighten-2">
          确定删除用户：{{ deleteItem.name }}？
        </v-card-title>
        <v-divider/>
        <v-card-actions>
          <v-btn color="primary" @click="deleteDialog=false">
            取消
          </v-btn>
          <v-spacer/>
          <v-btn color="warning" @click="confirmDelete">
            确定
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
        light
        v-model="addDialog"
        max-width="500"
    >
      <v-card>
        <v-card-title class="grey lighten-2">
          确定为{{ editedItem.name }}分配角色：{{rolesMap[editedRole][0]}}？
        </v-card-title>
        <v-divider/>
        <v-card-actions>
          <v-btn color="primary" @click="addDialog=false">
            取消
          </v-btn>
          <v-spacer/>
          <v-btn color="warning" @click="add(editedRole)">
            确定
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
        light
        v-model="revokeDialog"
        max-width="500"
    >
      <v-card>
        <v-card-title class="grey lighten-2">
          确定撤回{{ editedItem.name }}的角色：{{rolesMap[editedRole][0]}}？
        </v-card-title>
        <v-divider/>
        <v-card-actions>
          <v-btn color="primary" @click="revokeDialog=false">
            取消
          </v-btn>
          <v-spacer/>
          <v-btn color="warning" @click="revoke(editedRole)">
            确定
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
export default {
  name: "UsersTable",
  inject: ['message', 'reload'],
  data() {
    return {
      search: '',
      users: [],
      role: '',
      dialog: false,
      valid: false,
      password: '',
      rolesMap: {
        'sys_admin': ['系统管理员', 'indigo'],
        'teach_admin': ['教务老师', 'purple'],
        'assistant': ['助教', 'blue'],
        'teacher': ['教师', 'orange darken-3'],
        'student': ['学生', undefined]
      },
      roles: ['sys_admin', 'teach_admin', 'assistant', 'teacher'],
      chRoles: ['系统管理员', '教务老师', '助教', '教师','学生'],
      editedItem: {},
      deleteItem: {},
      deleteDialog: false,
      addDialog: false,
      revokeDialog: false,
      editedRole: 'teacher',
      unselected: [],
      passwordRules: [
        v => !!v || "密码不能为空！",
        v => (v && v.length >= 6) || "密码不能低于6位！",
      ],
    }
  },
  methods: {
    getUsers() {
      this.$axios.get('admin/users')
          .then(({data}) => {
            this.users = data.data
          })
          .catch(err => {
            console.log(err)
            this.message('error', '获取用户错误')
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
    openAddDialog(role){
      this.editedRole = role
      this.addDialog = true
    },
    openRevokeDialog(role){
      this.editedRole = role
      this.revokeDialog = true
    },
    add(role) {
      console.log(role)
      this.$axios.post(`admin/user/${this.editedItem.ID}/role`, {
        description: role,
      })
          .then(() => {
            this.message('info', '成功添加角色')
            this.reload()
          })
          .catch(() => {
            this.message('warning', '添加角色失败')
          })
    },
    revoke(role) {
      console.log(role)
      this.$axios.delete(`admin/user/${this.editedItem.ID}/role/${role}`)
          .then(() => {
            this.message('info', '成功撤回角色')
            this.reload()
          })
          .catch(() => {
            this.message('warning', '撤回角色失败')
          })
    },

    openDeleteDialog(item) {
      this.deleteItem = item
      this.deleteDialog = true
    },
    confirmDelete() {
      this.$axios.delete(`admin/user/${this.deleteItem.ID}`)
          .then(() => {
            this.message('info', `成功删除用户：${this.deleteItem.name}`)
            this.reload()
          })
          .catch(() => {
            this.message('error', `删除用户失败`)
          })
    },
    changePassword() {
      let data = new FormData()
      data.append('pwd', this.password)
      this.$axios.put(`admin/user/${this.editedItem.ID}`,data)
          .then(() => {
            this.message('info', `成功修改密码`)
            this.reload()
          })
          .catch(() => {
            this.message('error', `修改密码失败`)
          })
    }
  },
  computed: {
    headers() {
      return [
        {text: 'ID', value: 'ID'},
        {text: '姓名', value: 'name'},
        {text: '邮箱', value: 'email'},
        {text: '在线', value: 'online'},
        {
          text: '角色', value: 'roles',
          filter: value => {
            if (!this.role) return true
            return value.map(x => this.rolesMap[x.description][0]).includes(this.role)
          }
        },
        {text: '行为', value: 'actions', sortable: false,}
      ]
    },
  },
  created() {
    this.getUsers()
  },
}
</script>

<style scoped>

</style>
