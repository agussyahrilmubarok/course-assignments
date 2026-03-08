<script setup lang="js">
import Navbar from '@/components/Navbar.vue';
import Footer from '@/components/Footer.vue';
</script>

<template>
  <div class="project-page">
    <section class="dashboard-header pt-5">
      <div class="container mx-auto relative">
        <Navbar />
      </div>
    </section>
    <router-view></router-view>
    <div class="cta-clip"></div>
    <section class="call-to-action bg-[#3B41E3] pt-64 pb-10"></section>
    <Footer />
  </div>
</template>

<style lang="css" scoped>
.dashboard-header {
  background-image: url("@/assets/images/auth-background.svg");
  background-position: top right;
  background-repeat: no-repeat;
  background-color: #3b41e3;
  background-size: 250px;
  height: 100px;
}

.cta-clip {
  position: relative;
  top: 200px;
  bottom: 0;
  right: 0;
  left: 0;
  width: 100%;
  height: 300px;
  background-position: top right;
  background-size: 300px;
  background-repeat: no-repeat;
  background-color: #fff;
  transform: skewY(4deg);
  transform-origin: bottom right;
}

.call-to-action {
  background-image: url("@/assets/images/auth-background.svg");
  background-position: top right;
  background-repeat: no-repeat;
  background-size: 450px;
}
</style>
