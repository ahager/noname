<?php

    Class Flag {

        public static $log = [];

        public static function getClientConfig() {
            if (! isset($_COOKIE['flg-config'])) {
                return [];
            }
            return json_decode($_COOKIE['flg-config'], true);
        }

        public static function setClientConfig($json) {

            $config = self::getClientConfig();

            $config[$json['ClientId']]['ClientId'] = $json['ClientId'];
            $config[$json['ClientId']][$json['Name']]['Name'] = $json['Name'];
            $config[$json['ClientId']][$json['Name']]['Sticky'] = @$json['Sticky'];
            $config[$json['ClientId']][$json['Name']]['Status'] = $json['Status'];

            setcookie("flg-config", json_encode($config));

            return $config;
        }

        protected static function request($flagName, $forcedStatus='') {

            $clientId = '';
            if ($config = self::getClientConfig()) {
                $clientId = key($config);
            }

            if(isset($config[$clientId][$flagName])
                && isset($config[$clientId][$flagName]['Sticky']) && $config[$clientId][$flagName]['Sticky'] == 1) {
                $forcedStatus = $config[$clientId][$flagName]['Status'];
                $url = "http://localhost:8080/flag/$flagName/$clientId/$forcedStatus";
            } else {
                $url = "http://localhost:8080/flag/$flagName/$clientId";
            }

            self::$log[] = $url;

            $handle = curl_init();
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
            self::setClientConfig($json);
            return $json;
        }

        public static function isActive($flagName) {
            $result = self::request($flagName);
            return $result['Status'] == true;
        }

        public static function isInactive($flagName) {
            $result = self::request($flagName);
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
        header { height:100px; background:#0072BB; margin:1%; }
        main { float:left; height:400px; width:78%; background:#0072BB; margin:1%; }
        aside { float:right; height:400px; width:18%; background:#0072BB; margin:1%; }

        .inactive { background:#FF4C3B; }
    </style>
</head>
<body>
    <?php if(Flag::isActive('hero')): ?>
        <header>HERO</header>
    <?php else: ?>
        <header class="inactive">HERO</header>
    <?php endif; ?>

    <main>
        <p>MAIN CONTENT</p>
        <?php if(Flag::isActive('disruptor')): ?>
            <p>DISRUPTOR</p>
        <?php else: ?>
            <p class="inactive">DISRUPTOR</p>
        <?php endif; ?>
    </main>

    <?php if(Flag::isActive('sidebar')): ?>
        <aside>SIDEBAR</aside>
    <?php else: ?>
        <aside class="inactive">SIDEBAR</aside>
    <?php endif; ?>


    <div style="clear:both">
        <?php var_dump(Flag::getClientConfig()); ?>
        <hr>
        <?php var_dump(Flag::$log); ?>
    </div>
</body>
</html>
