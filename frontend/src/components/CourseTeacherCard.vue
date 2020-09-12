<template>
  <v-card>
    <v-card-title class="headline">教师</v-card-title>
    <v-card-actions>
      <v-row>
        <v-chip class="ml-2 my-1" v-for="(teacher,index) in course.teachers" :key="index">
          {{ teacher.user.name }}
        </v-chip>
      </v-row>
      <v-spacer/>
      <manage-teacher-dialog v-if="auth(['teach_admin','sys_admin'])"/>
    </v-card-actions>
  </v-card>
</template>

<script>
import ManageTeacherDialog from "@/components/ManageTeacherDialog";
export default {
  name: "CourseTeacherCard",
  components: {ManageTeacherDialog},
  props: ['course']
}
</script>

<style scoped>

</style>
