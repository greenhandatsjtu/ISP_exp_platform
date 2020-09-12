<template>
  <v-container fluid>
    <v-app-bar
        app
        color="grey"
        hide-on-scroll
    >
      <v-app-bar-nav-icon
          @click.stop="drawer.model = !drawer.model"
      />
      <div class="d-flex align-center">
        <v-img
            alt="SJTU Logo"
            class="shrink mr-2"
            contain
            :src="require('@/assets/logo.png')"
            transition="scale-transition"
            width="40"
        />
        <v-toolbar-title class="font-weight-black mr-2">
          信息安全综合实践平台
        </v-toolbar-title>
      </div>
    </v-app-bar>
    <v-navigation-drawer
        v-model="drawer.model"
        app
    >
      <!--            用户登录后才能看到用户名等信息-->
      <v-list v-if="$store.state.login">
        <v-list-item>
          <v-list-item-avatar>
            <v-img :src="require('@/assets/logo.png')"/>
          </v-list-item-avatar>
          <v-list-item-content>
            <v-list-item-title class="title font-weight-bold">{{ $store.state.user.name }}</v-list-item-title>
            <v-list-item-subtitle>{{ $store.state.user.email }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
        <v-row justify="space-around">
          <v-chip class="font-weight-bold" dark small v-for="(role,index) in $store.state.roles"
                  :key="index"
                  :color="rolesMap[role][1]">{{
              rolesMap[role][0]
            }}
          </v-chip>
        </v-row>
      </v-list>

      <v-divider></v-divider>
      <v-list dense nav rounded>
        <v-list-item-group color="primary" v-model="item">
          <v-list-item
              v-for="(item, index) in drawer.items"
              v-if="$store.state.login&&auth(item.roles)||item.to==='About'"
              :key="index"
              link
              :to="{name: item.to}"
          >
            <v-list-item-icon>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-icon>

            <v-list-item-content>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
      <template v-slot:append>
        <div class="px-2">
          <login-form v-if="!$store.state.login"/>
          <v-row v-else justify="space-around">
            <v-btn icon color="warning" @click="logout">
              <v-tooltip top>
                <template v-slot:activator="{on,attrs}">
                  <v-icon large v-on="on" v-bind="attrs">
                    mdi-logout
                  </v-icon>
                </template>
                <span>注销</span>
              </v-tooltip>
            </v-btn>
            <change-password-form/>
          </v-row>
          <v-switch
              v-model="$vuetify.theme.dark"
              label="dark theme"
          ></v-switch>
        </div>
      </template>
    </v-navigation-drawer>
  </v-container>
</template>

<script>
import LoginForm from "@/components/LoginForm";
import Roles from "@/components/Roles";
import ChangePasswordForm from "@/components/ChangePasswordForm";

export default {
  name: "MyAppBar",
  components: {ChangePasswordForm, LoginForm},
  data: () => ({
    item: 1,
    rolesMap: {
      'sys_admin': ['系统管理员', 'indigo'],
      'teach_admin': ['教务老师', 'purple'],
      'assistant': ['助教', 'blue'],
      'teacher': ['教师', 'orange darken-3'],
      'student': ['学生', undefined]
    },
    drawer: {
      model: null,
      items: [
        {title: '主页', icon: 'mdi-account-circle', to: 'Dashboard', roles: []},
        {title: '全部课程', icon: 'mdi-view-dashboard', to: 'AllCourses', roles: [Roles.sysAdmin, Roles.teachAdmin]},
        {title: '学生课程', icon: 'mdi-view-dashboard', to: 'StudentCourses', roles: [Roles.student]},
        {title: '教师课程', icon: 'mdi-view-dashboard', to: 'AdminCourses', roles: [Roles.teacher]},
        {title: '助教课程', icon: 'mdi-view-dashboard', to: 'AdminCourses', roles: [Roles.assistant]},
        {title: '用户管理', icon: 'mdi-account-multiple', to: 'Users', roles: [Roles.sysAdmin]},
        {title: '学生管理', icon: 'mdi-card-account-details', to: 'Students', roles: [Roles.sysAdmin, Roles.teachAdmin]},
        {title: '实验资源', icon: 'mdi-flask-empty-outline', to: 'AllResources', roles: [Roles.sysAdmin]},
        {title: '关于', icon: 'mdi-help-box', to: 'About', roles: []},
      ],
    },
  }),
  methods: {
    logout: function () {
      this.$axios.get('logout')
          .then(() => {
            this.$store.commit("logout")
            this.$router.push({name: "About"})
          })
    },
  }
}
</script>

<style scoped>

</style>
