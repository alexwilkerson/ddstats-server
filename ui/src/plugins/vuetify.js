import Vue from "vue";
import Vuetify from "vuetify/lib";

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    options: {
      customProperties: true
    },
    themes: {
      light: {
        background: "#f6f2ee",
        // background: "#fffefc",
        header: "#fffefc",
        // header: "#f6f2ee",
        footer: "#EBE7E4",
        primary: "#34302e",
        secondary: "#424242",
        accent: "#c33409",
        error: "#FF5252",
        info: "#2196F3",
        success: "#4CAF50",
        warning: "#FFC107"
      }
    }
  }
});
