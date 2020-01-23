<template>
  <div id="wrapper" v-if="loaded">
    <GameInfo :gameInfo="gameInfo"></GameInfo>
    <div class="chart-title">{{ titleGems }}</div>
    <div id="gems-chart" v-if="showGems">
      <apexchart
        type="area"
        height="200"
        width="800px"
        :options="chartOptionsGems"
        :series="seriesGems"
      ></apexchart>
    </div>
    <div class="chart-title">{{ titleHomingDaggers }}</div>
    <div id="homing-daggers-chart" v-if="showHomingDaggers">
      <apexchart
        type="area"
        height="200"
        width="800px"
        :options="chartOptionsHomingDaggers"
        :series="seriesHomingDaggers"
      ></apexchart>
    </div>
    <div class="chart-title">{{ titleAccuracy }}</div>
    <div id="accuracy-chart" v-if="showAccuracy">
      <apexchart
        type="area"
        height="200"
        width="800px"
        :options="chartOptionsAccuracy"
        :series="seriesAccuracy"
      ></apexchart>
    </div>
    <div class="chart-title">{{ titleEnemiesAlive }}</div>
    <div id="enemies-alive-chart" v-if="showEnemiesAlive">
      <apexchart
        type="area"
        height="200"
        width="800px"
        :options="chartOptionsEnemiesAlive"
        :series="seriesEnemiesAlive"
      ></apexchart>
    </div>
    <div class="chart-title">{{ titleEnemiesKilled }}</div>
    <div id="enemies-killed-chart" v-if="showEnemiesKilled">
      <apexchart
        type="area"
        height="200"
        width="800px"
        :options="chartOptionsEnemiesKilled"
        :series="seriesEnemiesKilled"
      ></apexchart>
    </div>
  </div>
</template>

<script>
import VueApexCharts from "vue-apexcharts";
import axios from "axios";
import GameInfo from "./GameInfo";
export default {
  data() {
    return {
      loaded: false,
      gameInfo: null,
      titleGems: "Gems: 0",
      titleHomingDaggers: "Homing Daggers: 0",
      titleAccuracy: "Accuracy: 0",
      titleEnemiesAlive: "Enemies Alive: 0",
      titleEnemiesKilled: "Enemies Killed: 0",
      showGems: false,
      showHomingDaggers: false,
      showAccuracy: false,
      showEnemiesAlive: false,
      showEnemiesKilled: false,
      chartOptionsGems: null,
      chartOptionsHomingDaggers: null,
      chartOptionsAccuracy: null,
      chartOptionsEnemiesAlive: null,
      chartOptionsEnemiesKilled: null,
      seriesGems: [
        {
          name: "Gems",
          data: []
        }
      ],
      seriesHomingDaggers: [
        {
          name: "Homing Daggers",
          data: []
        }
      ],
      seriesAccuracy: [
        {
          name: "Accuracy",
          data: []
        }
      ],
      seriesEnemiesAlive: [
        {
          name: "Enemies Alive",
          data: []
        }
      ],
      seriesEnemiesKilled: [
        {
          name: "Enemies Killed",
          data: []
        }
      ]
    };
  },
  components: {
    apexchart: VueApexCharts,
    GameInfo
  },
  methods: {
    getChartOptions: function(
      id,
      name,
      color,
      level2time,
      level3time,
      level4time
    ) {
      window.console.log(this);
      // let tooltipOffset = 0;
      // switch (name) {
      //   case "Gems":
      //     tooltipOffset = 60;
      //     break;
      //   case "Homing Daggers":
      //     tooltipOffset = 120;
      //     break;
      //   case "Accuracy":
      //     tooltipOffset = 60;
      //     break;
      //   case "Enemies Alive":
      //     tooltipOffset = 60;
      //     break;
      //   case "Enemies Killed":
      //     tooltipOffset = 60;
      //     break;
      // }
      let yaxisLabelFormatter = value => value;
      if (name === "Accuracy") yaxisLabelFormatter = value => value + "%";
      let annotations = [];
      if (level2time !== 0) {
        annotations.push({
          x: level2time,
          strokeDashArray: 5,
          borderColor: "#c33409",
          label: {
            borderColor: "transparent",
            style: {
              color: "#c33409",
              background: "transparent"
            },
            text: "L2",
            orientation: "horizontal",
            offsetY: -8
          }
        });
      }
      if (level3time !== 0) {
        annotations.push({
          x: level3time,
          strokeDashArray: 5,
          borderColor: "#c33409",
          label: {
            borderColor: "transparent",
            style: {
              color: "#c33409",
              background: "transparent"
            },
            text: "L3",
            orientation: "horizontal",
            offsetY: -8
          }
        });
      }
      if (level4time !== 0) {
        annotations.push({
          x: level4time,
          strokeDashArray: 5,
          borderColor: "#c33409",
          label: {
            borderColor: "transparent",
            style: {
              color: "#c33409",
              background: "transparent"
            },
            text: "L4",
            orientation: "horizontal",
            offsetY: -8
          }
        });
      }
      // let vm = this;
      return {
        chart: {
          id: id,
          group: "game",
          type: "area",
          height: 200,
          toolbar: {
            show: false
          },
          fontFamily:
            "'alte_haas_grotesk_bold', 'Helvetica Neue', Helvetica, Arial",
          animations: {
            enabled: false
          },
          zoom: {
            enabled: false
          }
          // events: {
          //   dataPointMouseEnter: function(event, chartContext, config) {
          //     window.console.log(event, chartContext, config);
          //   }
          // }
        },
        // title: {
        //   text: this.titleGems
        // },
        colors: [color],
        dataLabels: {
          enabled: false
        },
        fill: {
          colors: [color],
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
        annotations: {
          xaxis: annotations
        },
        markers: {
          strokeWidth: 2,
          strokeColors: "#c33409",
          hover: {
            size: 4
          }
        },
        xaxis: {
          type: "numeric",
          tickAmount: 16,
          crosshairs: {
            position: "front",
            stroke: {
              color: "#EBE7E4",
              dashArray: 5
            }
          },
          tooltip: {
            style: {
              fontFamily:
                "'alte_haas_grotesk_bold', 'Helvetica Neue', Helvetica, Arial"
            },
            formatter: value => value.toFixed(4) + "s"
          }
        },
        yaxis: {
          labels: {
            formatter: yaxisLabelFormatter,
            minWidth: 42
          },
          decimalsInFloat: false
        },
        tooltip: {
          marker: {
            show: false
          },
          custom: ({ series, seriesIndex, dataPointIndex }) => {
            switch (name) {
              case "Gems":
                this.titleGems = "Gems: " + series[seriesIndex][dataPointIndex];
                break;
              case "Homing Daggers":
                this.titleHomingDaggers =
                  "Homing Daggers: " + series[seriesIndex][dataPointIndex];
                break;
              case "Accuracy":
                this.titleAccuracy =
                  "Accuracy: " + series[seriesIndex][dataPointIndex] + "%";
                break;
              case "Enemies Alive":
                this.titleEnemiesAlive =
                  "Enemies Alive: " + series[seriesIndex][dataPointIndex];
                break;
              case "Enemies Killed":
                this.titleEnemiesKilled =
                  "Enemies Killed: " + series[seriesIndex][dataPointIndex];
                break;
            }
            return "";
            //   '<div class="tool-tip">' +
            //   "Homing Daggers: " +
            //   series[seriesIndex][dataPointIndex] +
            //   "</div>"
            // );
          }
        }
      };
    }
  },
  mounted() {
    axios
      .get(
        process.env.VUE_APP_API_URL +
          "/api/v2/game/full?id=" +
          this.$route.params.id
      )
      .then(response => {
        (this.chartOptionsGems = this.getChartOptions(
          "gems-chart",
          "Gems",
          colors[0],
          response.data.game_info.level_two_time,
          response.data.game_info.level_three_time,
          response.data.game_info.level_four_time
        )),
          (this.chartOptionsHomingDaggers = this.getChartOptions(
            "homing-daggers-chart",
            "Homing Daggers",
            colors[1],
            response.data.game_info.level_two_time,
            response.data.game_info.level_three_time,
            response.data.game_info.level_four_time
          ));
        this.chartOptionsAccuracy = this.getChartOptions(
          "accuracy-chart",
          "Accuracy",
          colors[2],
          response.data.game_info.level_two_time,
          response.data.game_info.level_three_time,
          response.data.game_info.level_four_time
        );
        this.chartOptionsEnemiesAlive = this.getChartOptions(
          "enemies-alive-chart",
          "Enemies Alive",
          colors[3],
          response.data.game_info.level_two_time,
          response.data.game_info.level_three_time,
          response.data.game_info.level_four_time
        );
        this.chartOptionsEnemiesKilled = this.getChartOptions(
          "enemies-killed-chart",
          "Enemies Killed",
          colors[4],
          response.data.game_info.level_two_time,
          response.data.game_info.level_three_time,
          response.data.game_info.level_four_time
        );
        for (let i = 0; i < response.data.states.length; i++) {
          if (response.data.states[i].gems > 0) {
            this.showGems = true;
          }
          if (response.data.states[i].homing_daggers > 0) {
            this.showHomingDaggers = true;
          }
          if (response.data.states[i].accuracy > 0) {
            this.showAccuracy = true;
          }
          if (response.data.states[i].enemies_alive > 0) {
            this.showEnemiesAlive = true;
          }
          if (response.data.states[i].enemies_killed > 0) {
            this.showEnemiesKilled = true;
          }
          this.seriesGems[0].data.push([
            response.data.states[i].game_time,
            response.data.states[i].gems
          ]);
          this.seriesHomingDaggers[0].data.push([
            response.data.states[i].game_time,
            response.data.states[i].homing_daggers
          ]);
          this.seriesAccuracy[0].data.push([
            response.data.states[i].game_time,
            response.data.states[i].accuracy
          ]);
          this.seriesEnemiesAlive[0].data.push([
            response.data.states[i].game_time,
            response.data.states[i].enemies_alive
          ]);
          this.seriesEnemiesKilled[0].data.push([
            response.data.states[i].game_time,
            response.data.states[i].enemies_killed
          ]);
        }
        this.gameInfo = response.data.game_info;
        this.loaded = true;
      })
      .catch(error => window.console.log(error));
  }
};

const colors = ["#898685", "#737271", "#5f5c5b", "#4a4746", "#373534"];
// const colors = ["#4a4746", "#4a4746", "#4a4746", "#4a4746", "#4a4746"];
</script>

<style scoped>
#wrapper {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  max-width: 800px;
  margin: 0 auto;
}
#gems-chart,
#homing-daggers-chart,
#accuracy-chart,
#enemies-alive-chart,
#enemies-killed-chart {
  width: 800px;
  margin: 0 auto;
}
.chart-title {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 15px;
  color: #34302e;
  padding-left: 6px;
}
</style>

<style>
.apexcharts-tooltip,
.apexcharts-tooltip.active,
.apexcharts-xaxistooltip,
.apexcharts-xaxistooltip.active,
.apexcharts-marker {
  transition: none !important;
}
.apexcharts-xaxistooltip {
  background: #fffefc;
  color: #c33409;
  border: #ebe7e4;
}
.apexcharts-xaxistooltip-bottom::before {
  border-bottom-color: #ebe7e4;
}
.apexcharts-xaxistooltip-bottom::after {
  border-bottom-color: #fffefc;
}
</style>
