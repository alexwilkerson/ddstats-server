<template>
  <div id="wrapper" v-if="loaded">
    <div id="gems-chart">
      <apexchart
        type="area"
        height="220"
        width="800px"
        :options="chartOptionsGems"
        :series="seriesGems"
      ></apexchart>
    </div>
    <div id="homing-daggers-chart"></div>
    <div id="accuracy-chart"></div>
    <div id="enemies-alive-chart"></div>
    <div id="enemies-killed-chart"></div>
  </div>
</template>

<script>
import VueApexCharts from "vue-apexcharts";
import axios from "axios";
export default {
  data() {
    return {
      loaded: false,
      chartOptionsGems: {
        chart: {
          id: "gems-chart",
          type: "area",
          height: 220,
          toolbar: {
            show: false
          },
          fontFamily:
            "'alte_haas_grotesk_bold', 'Helvetica Neue', Helvetica, Arial",
          animations: {
            enabled: false
          }
        },
        title: {
          text: "Gems"
        },
        dataLabels: {
          enabled: false
        },
        fill: {
          colors: ["#34302e"],
          opacity: 1,
          type: "solid"
        },
        stroke: {
          show: false
        },
        grid: {
          position: "front",
          borderColor: "#EBE7E4",
          strokeDashArray: 5,
          yaxis: {
            lines: {
              show: false
            }
          },
          xaxis: {
            lines: {
              show: true
            }
          }
        },
        xaxis: {
          type: "numeric",
          tickAmount: 16
        }
      },
      seriesGems: [
        {
          name: "Gems",
          data: []
        }
      ]
    };
  },
  components: {
    apexchart: VueApexCharts
  },
  mounted() {
    axios
      .get(process.env.VUE_APP_API_URL + "/api/v2/game/all?id=20")
      .then(response => {
        for (let i = 0; i < response.data.length; i++) {
          this.seriesGems[0].data.push([
            response.data[i].game_time,
            response.data[i].gems
          ]);
        }
        this.loaded = true;
      })
      .catch(error => window.console.log(error));
  }
};
</script>

<style scoped>
#wrapper {
  /* background: #fffefc; */
  margin: 0 auto;
  width: 100%;
}

#gems-chart {
  width: 800px;
  margin: 0 auto;
}
</style>
