<template>
  <v-card>
    <v-card-title class="headline">助教</v-card-title>
    <v-card-actions>
      <v-row>
        <v-chip class="ml-2 my-1" v-for="(assistant,index) in course.assistants" :key="index">
          {{ assistant.user.name }}
        </v-chip>
      </v-row>
      <v-spacer/>
      <manage-assistant-dialog v-if="auth(['teacher','sys_admin'])"/>
    </v-card-actions>
  </v-card>
</template>

<script>
import ManageAssistantDialog from "@/components/ManageAssistantDialog";
export default {
  name: "CourseAssistantCard",
  components: {ManageAssistantDialog},
  props: ['course'],
}
</script>

<style scoped>

</style>
