<template>
  <div id="chart">
    <h1>Deaths</h1>
    <apexchart type="bar" height="250px" :options="options" :series="series"></apexchart>
  </div>
</template>

<script>
import VueApexCharts from "vue-apexcharts";

export default {
  data() {
    return {
      options: {},
      series: [
        this.data.fallen,
        this.data.swarmed,
        this.data.impaled,
        this.data.gored,
        this.data.infested,
        this.data.opened,
        this.data.purged,
        this.data.desecrated,
        this.data.sacrificed,
        this.data.eviscerated,
        this.data.annihilated,
        this.data.intoxicated,
        this.data.envenmonated,
        this.data.incarnated,
        this.data.discarnated,
        this.data.barbed
      ]
    };
  },
  components: {
    apexchart: VueApexCharts
  },
  props: ["data"],
  mounted() {
    this.options = {
      plotOptions: {
        pie: {
          donut: {
            size: "60%"
          }
        }
      },
      labels: [
        "Fallen",
        "Swarmed",
        "Impaled",
        "Gored",
        "Infested",
        "Opened",
        "Purged",
        "Desecrated",
        "Sacrificed",
        "Eviscerated",
        "Annihilated",
        "Intoxicated",
        "Envenmonated",
        "Incarnated",
        "Discarnated",
        "Barbed"
      ],
      theme: {
        monochrome: {
          enabled: true,
          color: "#c33409"
        }
      }
    };
  }
};
</script>

<style>
</style>