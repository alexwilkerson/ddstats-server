<template>
  <div id="wrapper" v-if="loaded">
    <GameInfo :gameInfo="gameInfo"></GameInfo>
    <div class="chart-title" v-if="showGems">{{ titleGems }}</div>
    <div id="gems-chart" v-if="showGems">
      <apexchart
        ref="gems"
        type="area"
        height="200"
        :options="chartOptionsGems"
        :series="seriesGems"
      ></apexchart>
    </div>
    <div class="chart-title" v-if="showHomingDaggers">
      {{ titleHomingDaggers }}
    </div>
    <div id="homing-daggers-chart" v-if="showHomingDaggers">
      <apexchart
        ref="homing-daggers"
        type="area"
        height="200"
        :options="chartOptionsHomingDaggers"
        :series="seriesHomingDaggers"
      ></apexchart>
    </div>
    <div class="chart-title" v-if="showAccuracy">{{ titleAccuracy }}</div>
    <div id="accuracy-chart" v-if="showAccuracy">
      <apexchart
        ref="accuracy"
        type="area"
        height="200"
        :options="chartOptionsAccuracy"
        :series="seriesAccuracy"
      ></apexchart>
    </div>
    <div class="chart-title" v-if="showEnemiesAlive">
      {{ titleEnemiesAlive }}
    </div>
    <div id="enemies-alive-chart" v-if="showEnemiesAlive">
      <apexchart
        ref="enemies-alive"
        type="area"
        height="200"
        :options="chartOptionsEnemiesAlive"
        :series="seriesEnemiesAlive"
      ></apexchart>
    </div>
    <div class="chart-title" v-if="showEnemiesKilled">
      {{ titleEnemiesKilled }}
    </div>
    <div id="enemies-killed-chart" v-if="showEnemiesKilled">
      <apexchart
        ref="enemies-killed"
        type="area"
        height="200"
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
      homingDaggersMaxX: 0,
      homingDaggersMaxY: 0,
      enemiesAliveMaxX: 0,
      enemiesAliveMaxY: 0,
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
  computed: {
    theme() {
      return this.$vuetify.theme.dark ? "dark" : "light";
    }
  },
  watch: {
    theme: function() {
      this.resetChartOptions();
    }
  },
  components: {
    apexchart: VueApexCharts,
    GameInfo
  },
  methods: {
    resetChartOptions: function() {
      this.$refs.gems.updateOptions({
        grid: {
          borderColor: colors.grid[this.theme]
        },
        fill: {
          colors: [colors.area[this.theme]],
          opacity: 1,
          type: "solid"
        },
        xaxis: {
          labels: {
            style: {
              colors: colors.labels[this.theme]
            }
          }
        },
        yaxis: {
          labels: {
            style: {
              color: colors.labels[this.theme]
            }
          }
        }
      });
      // this.chartOptionsGems = { colors: ["#f00"] };
      // this.chartOptionsGems = this.getChartOptions(
      //   "gems-chart",
      //   "Gems",
      //   0,
      //   this.gameInfo.level_two_time,
      //   this.gameInfo.level_three_time,
      //   this.gameInfo.level_four_time
      // );
      // this.chartOptionsHomingDaggers = this.getChartOptions(
      //   "homing-daggers-chart",
      //   "Homing Daggers",
      //   1,
      //   this.gameInfo.level_two_time,
      //   this.gameInfo.level_three_time,
      //   this.gameInfo.level_four_time,
      //   this.gameInfo.homing_daggers_max,
      //   this.gameInfo.homing_daggers_max_time
      // );
      // this.chartOptionsAccuracy = this.getChartOptions(
      //   "accuracy-chart",
      //   "Accuracy",
      //   2,
      //   this.gameInfo.level_two_time,
      //   this.gameInfo.level_three_time,
      //   this.gameInfo.level_four_time
      // );
      // this.chartOptionsEnemiesAlive = this.getChartOptions(
      //   "enemies-alive-chart",
      //   "Enemies Alive",
      //   3,
      //   this.gameInfo.level_two_time,
      //   this.gameInfo.level_three_time,
      //   this.gameInfo.level_four_time,
      //   this.gameInfo.enemies_alive_max,
      //   this.gameInfo.enemies_alive_max_time
      // );
      // this.chartOptionsEnemiesKilled = this.getChartOptions(
      //   "enemies-killed-chart",
      //   "Enemies Killed",
      //   4,
      //   this.gameInfo.level_two_time,
      //   this.gameInfo.level_three_time,
      //   this.gameInfo.level_four_time
      // );
    },
    getChartOptions: function(
      id,
      name,
      i,
      level2time,
      level3time,
      level4time,
      max,
      maxTime
    ) {
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
      if (name === "Homing Daggers") {
        annotations.push({
          x: Math.ceil(maxTime),
          strokeDashArray: 5,
          borderColor: "#34302e",
          label: {
            borderColor: "transparent",
            style: {
              color: "#34302e",
              background: "#fffefc"
            },
            text: max.toString(),
            orientation: "horizontal",
            offsetY: -(5 - max.toString().length) * 3
          }
        });
      }
      if (name === "Enemies Alive") {
        annotations.push({
          x: Math.ceil(maxTime),
          strokeDashArray: 5,
          borderColor: "#34302e",
          label: {
            borderColor: "transparent",
            style: {
              color: "#34302e",
              background: "#fffefc"
            },
            text: max.toString(),
            orientation: "horizontal",
            offsetY: -(5 - max.toString().length) * 3
          }
        });
      }
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
        dataLabels: {
          enabled: false
        },
        fill: {
          colors: [colors.area[this.theme]],
          opacity: 1,
          type: "solid"
        },
        stroke: {
          show: false
        },
        grid: {
          position: "front",
          borderColor: colors.grid[this.theme],
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
          colors: colors.area[this.theme == "dark" ? "light" : "dark"],
          strokeWidth: 2,
          strokeColors: colors.area[this.theme],
          hover: {
            size: 4
          }
        },
        xaxis: {
          labels: {
            style: {
              colors: colors.labels[this.theme]
            }
          },
          type: "numeric",
          tickAmount: 16,
          crosshairs: {
            position: "front",
            stroke: {
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
            minWidth: 42,
            style: {
              color: colors.labels[this.theme]
            }
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
          0,
          response.data.game_info.level_two_time,
          response.data.game_info.level_three_time,
          response.data.game_info.level_four_time
        )),
          (this.chartOptionsHomingDaggers = this.getChartOptions(
            "homing-daggers-chart",
            "Homing Daggers",
            1,
            response.data.game_info.level_two_time,
            response.data.game_info.level_three_time,
            response.data.game_info.level_four_time,
            response.data.game_info.homing_daggers_max,
            response.data.game_info.homing_daggers_max_time
          ));
        this.chartOptionsAccuracy = this.getChartOptions(
          "accuracy-chart",
          "Accuracy",
          2,
          response.data.game_info.level_two_time,
          response.data.game_info.level_three_time,
          response.data.game_info.level_four_time
        );
        this.chartOptionsEnemiesAlive = this.getChartOptions(
          "enemies-alive-chart",
          "Enemies Alive",
          3,
          response.data.game_info.level_two_time,
          response.data.game_info.level_three_time,
          response.data.game_info.level_four_time,
          response.data.game_info.enemies_alive_max,
          response.data.game_info.enemies_alive_max_time
        );
        this.chartOptionsEnemiesKilled = this.getChartOptions(
          "enemies-killed-chart",
          "Enemies Killed",
          4,
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

const colors = {
  area: {
    light: "#4a4746",
    dark: "#E0DDDB"
  },
  grid: {
    light: "#EBE7E4",
    dark: "#1f1f1f"
  },
  labels: {
    light: "#34302e",
    dark: "#EBE7E4"
  }
};
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
  max-width: 800px;
  margin: 0 auto;
}
.chart-title {
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
  font-size: 15px;
  color: var(--v-primary-base);
  padding-left: 6px;
}

@media only screen and (max-width: 600px) {
  body {
    background-color: lightblue;
  }
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
