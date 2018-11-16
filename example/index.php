<?php

    Class Flag {

        public static $log = [];
        public static $config = [];

        public static function getClientConfig() {
            if (! isset($_COOKIE['flg-config'])) {
                return [];
            }
            return json_decode($_COOKIE['flg-config'], true);
        }

        public static function setClientConfig($json) {

            // self::$config = self::getClientConfig();

            // self::$log[] = json_encode($config);

            // $config[$json['ClientId']]['ClientId'] = $json['ClientId'];
            self::$config[$json['clientId']][$json['name']]['name'] = $json['name'];
            self::$config[$json['clientId']][$json['name']]['sticky'] = $json['sticky'];
            self::$config[$json['clientId']][$json['name']]['status'] = $json['status'];

            // self::$log[] = json_encode($config);

            setcookie("flg-config", json_encode(self::$config), time()+86400);

            return self::$config;
        }

        protected static function request($flagName, $forcedStatus='') {

            $clientId = '';
            if (! self::$config) {
                self::$config = self::getClientConfig();
            }

            if (self::$config) {
                $clientId = key(self::$config);
            }

            if(isset(self::$config[$clientId][$flagName])
                && isset(self::$config[$clientId][$flagName]['sticky']) && self::$config[$clientId][$flagName]['sticky'] == 1) {
                $forcedStatus = self::$config[$clientId][$flagName]['status'];
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
            return $result['status'] == true;
        }

        public static function isInactive($flagName) {
            $result = self::request($flagName);
            return $result['status'] == false;
        }
    }

?>
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>FLGR Example File</title>
    <style>
        header, main, aside { box-sizing:border-box; font-family: 'Helvetica Neue', Helvetica; font-size:48px; color:#fff; font-weight: 100; text-align:center; }
        header { height:100px; background: linear-gradient(135deg, #7abcff 0%,#60abf8 44%,#4096ee 100%); margin:1%; display: flex; align-items: center; justify-content: center;}
        main { float:left;  padding: 20px; height:400px; width:78%; background: linear-gradient(135deg, #7abcff 0%,#60abf8 44%,#4096ee 100%); margin:1%; }
        aside { float:right; height:400px; width:18%; background: linear-gradient(135deg, #7abcff 0%,#60abf8 44%,#4096ee 100%); margin:1%;  display: flex; align-items: center; justify-content: center;}

        .inactive { background: linear-gradient(135deg, #f22929 1%,#ff4949 100%); }
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
