<template>
  <div class="header-bar">
    <div class="sidebar-title">
      <el-icon class="menu-icon-btn" size="22" @click="foldClick" aria-label="切换菜单">
        <hamburger-button></hamburger-button>
      </el-icon>
      <nuxt-link class="title" to="/">{{ globalConfig.title }}</nuxt-link>
    </div>
    <div class="header-search">
      <div class="search-input">
        <input class="input" v-model="keywords" @keydown.enter="handleSearch" placeholder="搜索..." />
        <button class="btn" @click="handleSearch" :disabled="searchLoading" aria-label="搜索">
          <search-icon class="icon" size="16" />
        </button>
      </div>
    </div>
    <div v-if="!userLoading" class="header-right">
      <!-- 用户头像 -->
      <div v-if="isLoggedIn" class="avatar-box">
        <nuxt-link to="/space">
          <common-avatar :url="userInfo?.avatar" :size="40" :iconSize="22"></common-avatar>
        </nuxt-link>
        <div class="menu-container">
          <div class="transition"></div>
          <div class="header-menu">
            <div class="menu-info">
              <div class="name-box">
                <span class="name">{{ userInfo?.name }}</span>
                <span class="sign">{{ userInfo?.sign }}</span>
              </div>
            </div>
            <div class="divider disabled-select"></div>
            <nuxt-link class="menu-item disabled-select" to="/space">
              <div class="link-title">
                <user-icon class="icon"></user-icon>
                <span>个人中心</span>
              </div>
              <right-icon class="right-icon"></right-icon>
            </nuxt-link>
            <div class="menu-item disabled-select" @click="logout" :disabled="logoutLoading">
              <div class="link-title">
                <logout-icon class="icon"></logout-icon>
                <span>退出登录</span>
              </div>
              <right-icon class="right-icon"></right-icon>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="avatar-box">
        <div class="login-btn" @click="showLogin = true">登录</div>
      </div>
      <!-- 图形按钮 -->
      <nuxt-link class="icon-btn" to="/message/announce" aria-label="消息">
        <message-icon class="icon"></message-icon>
        <div class="icon-text">消息</div>
      </nuxt-link>
      <nuxt-link class="icon-btn" to="/space/history" aria-label="历史">
        <history-icon class="icon"></history-icon>
        <div class="icon-text">历史</div>
      </nuxt-link>
      <!-- 投稿按钮 -->
      <nuxt-link class="upload-btn disabled-select" to="/upload/video" aria-label="投稿">
        <upload-icon class="upload-icon"></upload-icon>
        <span>投稿</span>
      </nuxt-link>
    </div>
    <div v-else class="header-right"></div>
    <client-only>
      <login-dialog v-show="showLogin" @close="loginClose" @success="loginSuccess"></login-dialog>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { logoutAPI } from '@/api/auth';
import { getUserInfoAPI } from '@/api/user';
import LoginDialog from "@/components/login-dialog/index.vue";
import {
  HamburgerButton, Upload as UploadIcon, Search as SearchIcon,
  Message as MessageIcon, History as HistoryIcon,
  User as UserIcon, Logout as LogoutIcon, Right as RightIcon
} from '@icon-park/vue-next';

const emits = defineEmits(["changeFold"]);

const userLoading = ref(true);
const isLoggedIn = ref(false);
const showLogin = ref(false);
const userInfo = ref<UserInfoType>();
const searchLoading = ref(false);
const logoutLoading = ref(false);
const keywords = ref("");

// 获取用户信息
const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
    isLoggedIn.value = true;
  }
  userLoading.value = false;
}

// 左侧菜单折叠
const menuFold = ref(false);
const foldClick = () => {
  menuFold.value = !menuFold.value;
  emits("changeFold", menuFold.value);
}

// 退出登录
const logout = async () => {
  logoutLoading.value = true;
  try {
    await logoutAPI(storageData.get('refreshToken'));
    storageData.remove("token");
    storageData.remove('refreshToken');
    isLoggedIn.value = false;
  } finally {
    logoutLoading.value = false;
  }
}

// 搜索功能
const handleSearch = async () => {
  if (!keywords.value) {
    ElMessage.warning("请先输入搜索内容");
    return;
  }

  searchLoading.value = true;
  try {
    navigateTo(`/search/${keywords.value}`, {
      open: { target: '_blank' }
    });
  } finally {
    searchLoading.value = false;
  }
}

const loginClose = () => {
  showLogin.value = false;
}

// 登录成功
const loginSuccess = () => {
  getUserInfo();
  loginClose();
}

onBeforeMount(() => {
  getUserInfo();
})
</script>

<style lang="scss" scoped>
.header-bar {
  display: flex;
  align-items: center;
  width: 100%;
  height: 60px;
  box-shadow: inset 0 -1px #f1f2f3;

  .sidebar-title {
    display: flex;
    align-items: center;
    height: 100%;
    width: 220px;

    .menu-icon-btn {
      margin: 0 6px;
      padding: 6px;
      border-radius: 50%;
      cursor: pointer;

      &:hover {
        background-color: rgba(0, 0, 0, .1);
      }
    }

    .title {
      color: var(--font-primary-1);
      overflow: hidden;
      white-space: nowrap;
    }
  }

  .header-search {
    flex: 1;

    .search-input {
      width: 500px;
      position: relative;
      margin: 0 auto;

      .input {
        border: 1px solid #e5e5e5;
        outline: none;
        padding: 8px 30px 8px 11px;
        height: 36px;
        font-size: 12px;
        border-radius: 18px;
        width: 100%;
        color: var(--font-primary-1);
      }

      .btn {
        position: absolute;
        top: 0;
        right: 10px;
        border: none;
        width: 20px;
        height: 36px;
        font-size: 14px;
        background: transparent;
        cursor: pointer;

        .icon {
          display: block;
          width: 20px;
          color: var(--font-primary-5);
        }
      }
    }
  }

  .header-right {
    width: 286px;
    display: flex;
    align-items: center;

    .avatar-box {
      margin-right: 10px;
    }

    .login-btn {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      color: var(--primary-text-color);
      text-align: center;
      line-height: 40px;
      font-size: 14px;
      font-weight: 500;
      background-color: var(--primary-hover-color);
      cursor: pointer;
    }

    .upload-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: var(--primary-color);
      margin-left: 10px;
      width: 100px;
      height: 36px;
      border-radius: 5px;
      cursor: pointer;
      text-decoration: none;

      .upload-icon {
        width: 16px;
        height: 16px;
        margin-right: 5px;
      }

      &:hover {
        background-color: var(--primary-hover-color);
      }
    }
  }
}

.avatar-box {
  position: relative;
  cursor: pointer;

  &:hover .menu-container {
    display: block;
  }

  .menu-container {
    display: none;
    position: absolute;
    width: 230px;
    top: 40px;
    left: -110px;
    z-index: 999;
    transition: opacity 0.3s ease;

    .header-menu {
      padding: 12px;
      background-color: #fff;
      border-radius: 10px;
      box-shadow: 0 0 30px rgba(0, 0, 0, .1);

      .menu-info {
        display: flex;
        flex-direction: column;

        .name-box {
          padding-left: 12px;
          .name {
            font-size: 14px;
            color: var(--font-primary-2);
          }
          .sign {
            font-size: 12px;
            color: var(--font-primary-4);
            margin-top: 4px;
          }
        }
      }
    }
  }
}
</style>
