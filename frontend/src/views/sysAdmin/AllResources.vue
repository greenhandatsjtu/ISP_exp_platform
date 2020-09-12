<template>
  <v-container>
    <resources-table scope="all"/>
    <v-row justify="space-between">
      <v-card dark height="320" width="400" color="grey darken-2" class="mx-auto my-2" v-for="(metric, index) in metrics" :key="index">
        <v-card-title>{{metric.name}}</v-card-title>
        <v-row justify="center">
          <v-col cols="6" class="text-center">
            <v-list-item-title class="mb-2">CPU</v-list-item-title>
            <v-progress-circular
                :size="120"
                :width="10"
                :value="metric.cpu.percent*100"
                color="warning"
            >
              {{ (metric.cpu.percent*100).toFixed(2) }}%
            </v-progress-circular>
            <v-list-item>
              <v-list-item-title>可分配</v-list-item-title>
              {{metric.cpu.allocatable}}
            </v-list-item>
            <v-list-item>
              <v-list-item-title>已使用</v-list-item-title>
              {{metric.cpu.usage}}
            </v-list-item>
          </v-col>
          <v-col cols="6" class="text-center">
            <v-list-item-title class="mb-2">内存</v-list-item-title>
            <v-progress-circular
                :size="120"
                :width="10"
                :value="metric.memory.percent*100"
                color="info"
            >
              {{ (metric.memory.percent*100).toFixed(2) }}%
            </v-progress-circular>
            <v-list-item>
              <v-list-item-title>可分配</v-list-item-title>
              {{metric.memory.allocatable}}
            </v-list-item>
            <v-list-item>
              <v-list-item-title>已使用</v-list-item-title>
              {{metric.memory.usage}}
            </v-list-item>
          </v-col>
        </v-row>
      </v-card>
    </v-row>
  </v-container>
</template>

<script>
import ResourcesTable from "@/components/ResourcesTable";

export default {
  name: "AllResources",
  components: {ResourcesTable},
  data() {
    return {
      metrics: [],
      interval: {},
    }
  },
  methods: {
    getMetrics() {
      this.$axios.get('admin/resources/metrics')
          .then(({data}) => {
            this.metrics = data.data
          })
    }
  },
  beforeDestroy() {
    clearInterval(this.interval)
  },
  mounted() {
    this.interval = setInterval(() => {
     this.getMetrics()
    }, 20000)
  },
  created() {
    this.getMetrics()
  }
}
</script>

<style scoped>

</style>
