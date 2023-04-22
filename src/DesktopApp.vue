<template>
    <div></div>
</template>

<script>
export default {
    created() {
        if (process.env.NODE_ENV === 'production') {
            window.location.replace('../mobile/');
        } else {
            window.location.replace('../mobile.html');
        }
    }
}
</script>
