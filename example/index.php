<?php

    Class Flag {

        public static function clientId() {
            if (isset($_COOKIE['flg-id'])) {
                return $_COOKIE['flg-id'];
            }
            return '';
        }

        protected static function request($flagName, $clientId) {

            $handle = curl_init();
            $url = "http://localhost:8080/flag/$flagName/$clientId";
            curl_setopt($handle, CURLOPT_URL, $url);
            curl_setopt($handle, CURLOPT_USERAGENT, $_SERVER['HTTP_USER_AGENT']);
            curl_setopt($handle, CURLOPT_HTTPHEADER, [
                'Accept-Language: '.$_SERVER['HTTP_ACCEPT_LANGUAGE'],
                'Remote-Address: '.$_SERVER['REMOTE_ADDR']
            ]);
            curl_setopt($handle, CURLOPT_RETURNTRANSFER, true);
            $result = curl_exec($handle);
            curl_close($handle);

            $json = json_decode($result, true);
            if (isset($json['clientId'])) {
                setcookie("flg-id", $json['clientId']);
            }

            return $json;
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
    <style>
        header, main, aside { box-sizing:border-box; font-family:Helvetica; font-size:48px; color:#fff; font-weight:bold; text-align:center; }
        header { height:100px; background:#FF4C3B; margin:1%; }
        main { float:left; height:400px; width:78%; background:#0072BB; margin:1%; }
        aside { float:right; height:400px; width:18%; background:#FF4C3B; margin:1%; }
    </style>
</head>
<body>
    <?php if(Flag::isActive('my-hero')): ?>
        <header>HERO</header>
    <?php endif; ?>

    <main>MAIN CONTENT</main>

    <?php if(Flag::isActive('my-sidebar')): ?>
        <aside>SIDEBAR</aside>
    <?php endif; ?>
</body>
</html>
