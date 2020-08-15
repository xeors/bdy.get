<?
    function Http($url,$data=null,$header=null,$rh=0,$nb=0){
        $ch=curl_init();
        curl_setopt($ch,CURLOPT_URL,$url);
        curl_setopt($ch,CURLOPT_HEADER,$rh);
        curl_setopt($ch,CURLOPT_NOBODY,$nb);
        curl_setopt($ch,CURLOPT_RETURNTRANSFER,1);
        $header==null?:curl_setopt($ch,CURLOPT_HTTPHEADER,$header);
        $data==null?:(curl_setopt($ch,CURLOPT_POST,1)&&curl_setopt($ch,CURLOPT_POSTFIELDS,$data));
        $rdata=curl_exec($ch);
        curl_close($ch);
        return $rdata;
    }
    $path='bdy.inf';
    $appname='bdy.get';
    $version='1.0.0p';
    if(isset($_GET['h'])){
        echo $appname.'/'.$version.
        '<br>参数：<br>
        -B BDUSS<br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Cookie:BDUSS：在baidu.com下<br>
        -S STOKEN<br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Cookie:STOKEN：在pan.baidu.com下<br>
        -c	创建或修改BDUSS和STOKEN<br>
        -d ShareID<br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;ShareID：分享ID<br>
        -h	显示帮助<br>
        -p Password<br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Password：提取码<br>
        -t	测试BDUSS和STOKEN<br>
        -v	显示版本<br>
      参数优先级h>v>c>t。下载文件需要设置User-Agent:LogStatistic<br>
        ';exit;
    }else if(isset($_GET['v'])){
        echo $appname.'/'.$version;exit;
    }
    if((isset($_GET['B'])&&$_GET['B']!='')&&(isset($_GET['S'])&&$_GET['S']!='')){
        $STOKEN=$_GET['S'];
        $BDUSS=$_GET['B'];
    }else if((isset($_POST['B'])&&$_POST['B']!='')&&(isset($_POST['S'])&&$_POST['S']!='')){
        $BDUSS=$_POST['B'];
        $STOKEN=$_POST['S'];
        $r=file_put_contents($path,$BDUSS.';;'.$STOKEN);
        if($r==false){
            echo '写入失败!!';
            exit;
        }
    }else{
        if(file_exists($path)){
            $s=file_get_contents($path);
            $s=explode(';;',$s);
            if(count($s)==2){
                $BDUSS=$s[0];
                $STOKEN=$s[1];
            }
        }else{
            echo '<span>首次使用需要配置信息，请输入以下信息：</span>
            <form action="" method="POST">
            <span>BDUSS:</span><input type="text" name="B" placeholder="BDUSS"><br>
            <span>STOKEN:</span><input type="text" name="S" placeholder="STOKEN"><br>
            <input type="submit" value="Set."></form>';
            exit;
        }
    }
    if(isset($_GET['c'])&&!isset($_POST['B'])){
        echo '<span>创建或修改以下信息：</span>
        <form action="" method="POST">
        <span>BDUSS:</span><input type="text" name="B" placeholder="BDUSS"><br>
        <span>STOKEN:</span><input type="text" name="S" placeholder="STOKEN"><br>
        <input type="submit" value="Set."></form>';
        exit;
    }
    if(isset($_GET['t'])){
        if(!isset($_GET['B'])||$_GET['B']==''||!isset($_GET['S'])||$_GET['S']==''){
            if(file_exists($path)){
                $s=file_get_contents($path);
                $s=explode(';;',$s);
                if(count($s)==2){
                    $BDUSS=$s[0];
                    $STOKEN=$s[1];
                }
            }
        }
        $h=array(
            'Cookie:STOKEN='.$STOKEN.';BDUSS='.$BDUSS
        );
        $u='https://pan.baidu.com/disk/home';
        $r=Http($u,null,$h);
        if(mb_strlen($r)>23){
            echo '登陆信息有效!!';exit;
        }else{
            echo '登陆信息无效!!';exit;
        }
    }
    if((isset($_GET['d'],$_GET['p'])&&$_GET['d']!=''&&$_GET['p']!='')){
        $sl=$_GET['d'];
        if(preg_match('/baidu.com\/s\/(.*?) /',$sl,$sd)&&preg_match('/码: (.*?)$/',$sl,$pw)){
            $sid=trim($sd[1]);
            $pwd=trim($pw[1]);
        }else{
            $sid=$_GET['d'];
            $pwd=$_GET['p'];
        }
        if(substr($sid,0,1)==1){
            $sid=substr($sid,1,strlen($sid));
        }
        $url1='https://pan.baidu.com/share/verify?surl='.$sid;
        $data1='pwd='.$pwd;
        $header1=array(
            'User-Agent:netdisk',
            'Referer:https://pan.baidu.com/disk/home'
        );
        $r1=json_decode(Http($url1,$data1,$header1),1);
        if($r1['errno']==0){
            $randsk=$r1['randsk'];
            $url2='https://pan.baidu.com/s/1'.$sid;
            $header2=array(
                'Cookie:STOKEN='.$STOKEN.';BDUSS='.$BDUSS.';BDCLND='.$randsk,
                'User-Agent:Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.514.1919.810 Safari/537.36'
            );
            $r2=Http($url2,null,$header2);
            if(preg_match('/yunData.setData\(({.+)\);/',$r2,$m2)){
                $s=json_decode($m2[1],1);
                $sign=$s['sign'];$timestamp=$s['timestamp'];$shareid=$s['shareid'];$uk=$s['uk'];
                $url3='https://pan.baidu.com/share/list?app_id=250528&channel=chunlei&clienttype=0&desc=0&num=100&order=name&page=1&root=1&shareid='.$shareid.'&showempty=0&uk='.$uk.'&web=1';
                $r3=json_decode(Http($url3,null,$header2),1);
                if($r3['errno']==0){
                    $fsid=$r3['list'][0]['fs_id'];
                    $url4='https://pan.baidu.com/api/sharedownload?app_id=250528&channel=chunlei&clienttype=12&sign='.$sign.'&timestamp='.$timestamp.'&web=1';
                    $data4='encrypt=0&extra={"sekey":"'.$randsk.'"}&fid_list=['.$fsid.']&primaryid='.$shareid.'&uk='.$uk.'&product=share&type=nolimit';
                    $r4=json_decode(Http($url4,$data4,$header2),1);
                    if($r4['errno']==0){
                        $url5=$r4['list'][0]['dlink'];
                        $header5=array(
                            'User-Agent:LogStatistic',
                            'Cookie:BDUSS='.$BDUSS
                        );
                        $r5=Http($url5,null,$header5,1,0);
                        if(preg_match('/Location:(.*?)\\r\\n/',$r5,$m5)){
                            if(isset($_GET['r'])){
                                header('Location:'.$m5[1]);
                            }else{
                                echo trim($m5[1]);
                            }
                        }else{
                            echo '提取码错误或文件失效或登陆失效!!';
                        }
                    }                    
                }
            }
        }
    }else{
        echo '<span>输入以下信息来获取直链：</span>
        <form action="" method="GET">
        <span>ShareID:</span><input type="text" name="d" placeholder="分享ID，或者分享链接"><br>
        <span>Password:</span><input type="text" name="p" placeholder="提取码"><br>
        <input type="submit" value="Get."></form>';
        exit;
    }
?>