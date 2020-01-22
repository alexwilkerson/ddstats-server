<template>
  <p>{{ info }}</p>
</template>

<script>
import axios from "axios";
export default {
  data() {
    return {
      info: null
    };
  },
  mounted() {
    axios
      .get(process.env.VUE_APP_API_URL + "/api/v2/game/all?id=2")
      .then(response => (this.info = response))
      .catch(error => console.log(error));
  }
};
</script>

<style></style>
