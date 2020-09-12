<template>
  <v-card light>
    <v-card-title class="font-weight-bold headline grey lighten-1">YAML</v-card-title>
    <v-divider/>
    <v-card-subtitle v-if="!experiment.upload">暂无</v-card-subtitle>
    <v-list v-else dense>
      <v-list-item v-for="(file, index) in filesName" :key="index">
        <v-list-item-content>
          <v-list-item-title>{{ file }}</v-list-item-title>
        </v-list-item-content>
        <v-list-item-icon>
          <v-btn icon @click="download(file)">
            <v-icon color="success">mdi-download</v-icon>
          </v-btn>
        </v-list-item-icon>
      </v-list-item>
    </v-list>
    <v-card-actions>
      <upload-yaml-dialog/>
    </v-card-actions>
  </v-card>
</template>

<script>
import UploadYamlDialog from "@/components/UploadYamlDialog";
export default {
  name: "YamlCard",
  components: {UploadYamlDialog},
  props: ['experiment'],
  computed: {
    filesName() {
      return [`${this.$route.params.id}.yaml`, `${this.$route.params.id}_svc.yaml`]
    },
  },
  methods: {
    download(file){
      window.open(`api/admin/experiment/${this.$route.params.id}/yaml/${file}`)
    }
  },
}
</script>

<style scoped>

</style>
