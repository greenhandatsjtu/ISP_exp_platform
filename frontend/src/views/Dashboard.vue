<template>
  <v-container>
    <v-card max-width="400">
      <v-card-title class="font-weight-bold headline teal white--text">
        通知
      </v-card-title>
      <v-card v-for="(notice, index) in notices" :key="index" class="mx-2 my-3">
        <v-card-subtitle class="text--primary grey lighten-1">{{ notice.title }}</v-card-subtitle>
        <v-card-text>
         <v-row>
           <p class="text--secondary">{{ notice.author }}</p>
           <v-spacer/>
           <p class="text--secondary">{{new Date(notice.UpdatedAt).toLocaleString()}}</p>
         </v-row>
          <v-row>
            {{ notice.body }}
          </v-row>
        </v-card-text>
      </v-card>
    </v-card>
    <AddNoticeCard v-if="auth(['teach_admin','sys_admin','teacher','assistant'])"></AddNoticeCard>
  </v-container>
</template>

<script>
import AddNoticeCard from "../components/AddNoticeCard";


export default {
  name: "Dashboard",
  components:{AddNoticeCard},
  inject: ['message','reload'],
  data() {
    return {
      notices: [],
      flag: false,
      sheet: false,
    }
  },
  methods:{
  },
  beforeCreate() {
    this.$axios.get('notice')
        .then(({data}) => {
          this.notices = data.data
        })
  },
}
</script>

<style scoped>

</style>
