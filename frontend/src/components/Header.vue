<template>
    <n-space justify="space-between" align="center" style="height: 100%;">
      <div>
        <n-space>
          <n-avatar size="medium" color="#18a058">
            LOGO
          </n-avatar>
          <span style="font-size: 18px; font-weight: bold;">网站标题</span>
        </n-space>
      </div>
      <n-space>
        <RouterLink to="/" style="margin-right: 16px;">首页</RouterLink>
        <RouterLink to="/about" style="margin-right: 16px;">关于我们</RouterLink>
        <RouterLink to="/contact">联系我们</RouterLink>
      </n-space>
    </n-space>
  </template>
  