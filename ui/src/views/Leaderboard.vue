<template>
  <div class="wrapper">
    <LeaderboardTable :data="data" />
  </div>
</template>

<script>
import LeaderboardTable from "../components/LeaderboardTable";
import axios from "axios";

export default {
  data: () => ({
    loading: true,
    data: null
  }),
  components: {
    LeaderboardTable
  },
  mounted() {
    axios
      .get(
        process.env.VUE_APP_API_URL +
          "/api/v2/game/full?id=" +
          this.$route.params.id
      )
      .then(response => {});
  }
};
</script>

<style scoped>
.wrapper {
  padding-top: 40px;
  padding-bottom: 40px;
  max-width: 800px;
  margin: 0 auto;
}
</style>
