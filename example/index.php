<?php

    Class Flag {

        public static function clientId() {
            if (isset($_COOKIE['flg'])) {
                return $_COOKIE['flg'];
            } else {
                $clientId = md5($_SERVER['HTTP_USER_AGENT'].$_SERVER['REMOTE_ADDR']);
                setcookie("flg", $clientId);
                return $clientId;
            }
        }

        protected static function request($flagName, $clientId) {
            $handle = curl_init();
            $url = "http://localhost:8080/flag/$flagName/$clientId";
            curl_setopt($handle, CURLOPT_URL, $url);
            curl_setopt($handle, CURLOPT_RETURNTRANSFER, true);
            $result = curl_exec($handle);
            // var_dump($url, $result);
            curl_close($handle);
            return json_decode($result, true);
        }

        public static function isActive($flagName) {
            $result = self::request($flagName, self::clientId());
            return $result['Status'] == true;
        }

        public static function isInactive($flagName) {
            $result = self::request($flagName, self::clientId());
            return $result['Status'] == false;
        }
    }

?>
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>FLGR Example File</title>
</head>
<body>
    <?php if(Flag::isActive('feature-one')): ?>
        <section>Feature One is active!</section>
    <?php else: ?>
        <section>No Feature One!</section>
    <?php endif; ?>

</body>
</html>
