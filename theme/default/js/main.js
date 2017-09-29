var delay=1;
var tmr,tmn,tmg;
var THB=0;
var currency='THB';
var curr='thb';
var data;
var loaded=false;
var pass='';
var coin='';
var market={'bx':'bx.in.th'};
var last_coin='';
function network(s){
  var t=$('.href_'+s.replace('/',''));
  if(t.length>0) {
    pair=t.data('pair');
    $('#pairing').html('Loading...');
    $('#waiting').css('display','block');
    $('#ui').css('display','none');
    loaded=false;
    $('#navbarNetwork .active').removeClass('active');
    t.addClass('active');
    if(tmg)
    {
      clearInterval(tmg);
      tmg='';
    }
  }
}
function sendorder(t){
  pass=(t.save.value?t.password.value:'');
  $('#order_modal').modal('hide');
  var act=t.action.value,q='pair='+pair+'&action='+act+'&pass='+encodeURIComponent(t.password.value);
  if(act=='cancel'){
    q+='&id='+t.id.value;
  }else if(act=='order'){
    q+='&type='+t.type.value;
    q+='&rate='+t.rate.value;
    q+='&amount='+t.amount.value;
  }
  t.reset();
  $.getJSON(
    '/ajax/order?'+q,
    function(da){
      $(da.target).popover({
        title:(da.success?'สำเร็จ':'ผิดพลาด'),
        content:da.message,
      }).popover('show');
      setTimeout(function(){$(da.target).popover('dispose');},5000);
    }
  );
}
function _getData(){
  jQuery.ajax({
    type: "GET",
    url: '/ajax/data?pair='+pair+'&callback=?',
    dataType: "jsonp",
    success: function(da){
      var d = new Date,v;
      data = da;
      currency=data.pair.primary_currency;
      $('#time').html(d.getDate()+'/'+('0'+(d.getMonth()+1)).substr(-2)+'/'+d.getFullYear()+' '+('0'+d.getHours()).substr(-2)+':'+('0'+d.getMinutes()).substr(-2)+':'+('0'+d.getSeconds()).substr(-2));
      THB = data.usdthb.rates.THB;
      $('.usdthb').html('1 USD = '+THB+' THB');
      var tmp='';
      for(var i=0;i<data.sort.length;i++){
        var k,v;
        if((k=data.sort[i])&&(v=data.list[k])&&(data.bitfinex[v.secondary])&&(data.bittrex[v.secondary])){
          var bth=data.bitfinex[v.secondary][6];
          var btt=data.bittrex[v.secondary].Price;
          var p=(((bth*THB)-v.price)/v.price)*100;
          var p2=(((btt*THB)-v.price)/v.price)*100;
          var nv=$('#nav-'+v.primary+'-'+v.secondary);
          if(nv.length>0){
            nv.find('.'+v.primary+'_'+v.secondary).html(_num(_fs(v.price)));
            nv.find('.THB2_'+v.secondary).html(_num(_fs(bth*THB)));
            nv.find('.bth-change').attr('class',p>0?'green':'red').html(_num(_fs(p)))
            nv.find('.THB3_'+v.secondary).html(_num(_fs(btt*THB)));
            nv.find('.btt-change').attr('class',p2>0?'green':'red').html(_num(_fs(p2)))
            nv.find('.USD_'+v.secondary).html(_num(_fs(bth)));
            nv.find('.USD2_'+v.secondary).html(_num(_fs(btt)));
          }else{
            var ps='(<span class="bth-change '+(p>0?'green">+':'red">')+_num(_fs(p))+'</span>%)';
            var ps2='(<span class="btt-change '+(p2>0?'green">+':'red">')+_num(_fs(p2))+'</span>%)';
            tmp='<li class="nav-item" id="nav-'+v.primary+'-'+v.secondary+'"><a class="nav-link href_'+v.secondary+(pair==k?' active':'')+'" href="/'+v.secondary+'" data-pair="'+k+'">'+
            v.secondary+'<br>'+
            v.primary+': <span class="'+v.primary+'_'+v.secondary+'">'+_num(_fs(v.price))+'</span><br>'+
            'THB: <span class="THB2_'+v.secondary+'">'+_num(_fs(bth*THB))+'</span> '+ps+'<br>'+
            'THB: <span class="THB3_'+v.secondary+'">'+_num(_fs(btt*THB))+'</span> '+ps2+'<hr>'+
            'USD: <span class="USD_'+v.secondary+'">'+_num(_fs(bth))+'</span><br>'+
            'USD: <span class="USD2_'+v.secondary+'">'+_num(_fs(btt))+'</span>'+
            '</a></li>';
            $('#pair').append(tmp);
          }
        }
      };
      //$('#pair').html(tmp);
      $.each(data.wallet,function(k,v){
        if(v.total)
        {
          var wl=$('#wallet_'+k);
          if(!wl.length)
          {
            $('#wallet').append('<div id="wallet_'+k+'" class="col-6 col-sm-4 col-md-3 col-lg-2"></div>');
            wl=$('#wallet_'+k);
          }
          wl.html(k+': '+(v.available?_num(v.available):0)+' <em>(<span>'+_num(v.total)+'</span>)</em>');
        }
      });
      coin=data.pair.secondary_currency;
      if(!tmg)
      {
        tmg=setInterval(_getGraph,5000)
        _getGraph();
      }
      var cur=currency+'/'+data.pair.secondary_currency;
      document.title = _num(_fs(data.pair.last_price)) + ' : '+currency+'/'+data.pair.secondary_currency+' - BXGo v. '+ver;
      $('#pairing').html(data.pair.secondary_currency);
      $('.'+currency+'_'+data.pair.secondary_currency).html(_num(_fs(data.pair.last_price)));
      $('.primary_currency').html(currency);
      $('.secondary_currency').html(data.pair.secondary_currency);
      $('#bx_price').html(_num(_fs(data.pair.last_price)));
      _color($('#bx_change'),data.pair.change);
      $('#bx_vol').html(_num(data.pair.volume_24hours));
      $('#bx_buy').html(_num(_fs(data.pair.orderbook.bids.highbid)));
      $('#bx_buy_vol').html(_num(data.pair.orderbook.bids.volume));
      $('#bx_sell').html(_num(_fs(data.pair.orderbook.asks.highbid)));
      $('#bx_sell_vol').html(_num(data.pair.orderbook.asks.volume));

      var bfn = data.bitfinex[data.pair.secondary_currency];
      if(bfn){
        $('#bfn_price').html(_num(_fs(bfn[6]*THB)));
        $('#bfn_price_usd').html(_num(bfn[6]));
        _color($('#bfn_change'),_fs(bfn[5]*100));
        $('#bfn_vol').html(_num(bfn[7]));
        $('#bfn_buy').html(_num(_fs(bfn[0]*THB)));
        $('#bfn_buy_usd').html(_num(_fs(bfn[0])));
        $('#bfn_buy_vol').html(_num(bfn[1]));
        $('#bfn_sell').html(_num(_fs(bfn[2]*THB)));
        $('#bfn_sell_usd').html(_num(_fs(bfn[2])));
        $('#bfn_sell_vol').html(_num(bfn[3]));
      }
      var btx = data.bittrex[data.pair.secondary_currency];
      if(btx){
        $('#btx_price').html(_num(_fs(btx.Price*THB)));
        $('#btx_price_usd').html(_num(_fs(btx.Price)));
        _color($('#btx_change'),_fs(btx.Change*100));
        $('#btx_vol').html(_num(btx.Volume));
        $('#btx_buy').html(_num(_fs(btx.Bid*THB)));
        $('#btx_buy_usd').html(_num(_fs(btx.Bid)));
        $('#btx_buy_vol').html('Order: '+_num(btx.Order_Buy));
        $('#btx_sell').html(_num(_fs(btx.Ask*THB)));
        $('#btx_sell_usd').html(_num(_fs(btx.Ask)));
        $('#btx_sell_vol').html('Order: '+_num(btx.Order_Sell));
      }
      $('#trend_trade').html(_num(_fs(data.trend.TRADE_AVG))).attr('class',(data.trend.TRADE_AVG<data.pair.last_price?'bred':'bgreen'));
      $('#trend_10_price').html(_num(data.trend.Price_AVG_10));
      _color($('#trend_10_change'),_fs(data.trend.Price_AVG_10-data.pair.last_price));
      $('#trend_10_bid_sell').html(_num(data.trend.UP_AVG_10));
      $('#trend_10_bid_buy').html(_num(data.trend.DOWN_AVG_10));
      $('#trend_10_vol_sell').html(_num(data.trend.UP_Vol_10));
      $('#trend_10_vol_buy').html(_num(data.trend.DOWN_Vol_10));

      var tmp='',si;
      $.each(data.order,function(k,v){
        si=_sim(v.order_type,v.rate,v.amount)
        tmp+='<tr class="'+v.order_type+'"><td>'+v.order_type+'</td><td>'+_num(_fs(v.rate))+' <em>('+(v.order_type=='sell'?'<span class="red2">+'+_num(_fs(v.rate-data.pair.last_price))+'</span>':'<span class="green2">'+_num(_fs(v.rate-data.pair.last_price))+'</span>')+')</em></td><td>'+_num(v.amount)+'</td><td>'+v.date+'</td><td>'+si.do+'</td><td>'+si.profit+'</td><td width="20px"><button type="button" class="btn btn-bxgo" data-toggle="modal" data-target="#order_modal" data-action="cancel" data-id="'+v.order_id+'"><i class="fa fa-trash" aria-hidden="true"></i></button></td></tr>';
      });
      $('#order').html(tmp);
      $('#order_count').html('buy: '+$('#order>tr.buy').length+' / sell: '+$('#order>tr.sell').length);

      var tmp='';
      $.each(data.trans,function(k,v){
        var c = (v.Primary>0?'sell':'buy');
        var rate=Math.abs(v.Primary/v.Secondary);
        tmp+='<tr class="'+c+'"><td>'+(v.Primary>0?'+':'')+_num(_fs(v.Primary))+'</td><td>'+(v.Secondary>0?'+':'')+_num(v.Secondary)+'</td><td>'+_num(_fs(rate))+'</td><td>'+_num(v.Fee)+'</td><td>'+v.Date+'</td></tr>';
      });
      $('#history').html(tmp);
      $('#history_count').html('buy: '+$('#history>tr.buy').length+' / sell: '+$('#history>tr.sell').length);

      $('#conf_enable').attr('class',(data.conf.Enable?'green':'red')).html(data.conf.Enable?'Enabled':'Disabled')
      $('#conf_buy_budget').html(_num(_fs(data.conf.Budget)));
      $('#conf_current_budget').html(_num(_fs(data.wallet&&data.wallet[currency]?data.wallet[currency].available:'-')));//+' '+currency);

      $('#conf_buy_max_price').html(_num(_num(_fs(data.conf.Max_Price))));
      $('#conf_current_price').html(_num(_fs(data.pair.last_price)));//+' '+currency);
      $('#conf_buy_max_order').html(_num(data.conf.Max_Order));
      $('#conf_current_order').html($('#order>tr').length);
      $('#conf_diff').html(_num(data.conf.Cycle));
      $('#conf_diff_val').html(_num(_fs((data.conf.Cycle/100)*data.pair.last_price)));//+' '+currency);
      $('#conf_sell_margin').html(_num(data.conf.Margin));
      $('#conf_sell_margin_val').html(_num(_fs((data.conf.Margin/100)*data.pair.last_price)))//+' '+currency);

      var tmp='';
      $.each(data.sims,function(k,v){
        tmp+='<tr><td class="'+(v.Order_Buy?'in':'out')+'order">'+_num(_fs(v.Buy))+' '+currency+(v.Order_Buy?' (order:'+_num(_fs(v.Order_Buy))+')':'')+'</td><td class="'+(v.Order_Sell?'in':'out')+'order">'+_num(_fs(v.Sell))+' '+currency+(v.Order_Sell?' (order:'+_num(_fs(v.Order_Sell))+')':'')+'</td><td>'+_num(_fs(v.Margin))+' '+currency+'</td><td>'+_num(v.Coin.toFixed(8))+' '+data.pair.secondary_currency+'</td><td>'+_num(_fs(v.Profit))+' '+currency+'</td><td>'+_num(_fs(v.Diff))+' '+currency+'</td></tr>';
      });
      $('#sims').html(tmp);

      var tmp=[];
      if(data.delay.Next_Buy)
      {
        tmp.push('buy: '+data.delay.Next_Buy);
      }
      if(data.delay.Next_Sell)
      {
        tmp.push('sell: '+data.delay.Next_Sell);
      }
      $('#delay').html((tmp.length>0?'Delay for next - '+tmp.join(', '):''));
      if(!loaded)
      {
        loaded=true;
        $('#waiting').css('display','none');
//        $('#bxgo').css('display','block');
        $('#ui').css('display','block');
        $('#mywall').css('display','block');
      }
    },
    error: function(XMLHttpRequest, textStatus, errorThrown){

    }
  });
}
function _fs(s)
{
  if(s && !isNaN(s) && s.toString().indexOf('.')>0)
  {
    if(s.toFixed)s=s.toFixed(2);
    var t=s.toString().split('.');
    return t[0]+'.'+(t[1]?t[1]+'0':'00').substr(0,2);
  }else if(s && !isNaN(s)){
    return s+'.00';
  }
  return s?s:'0.00';
}
function _color(o,i)
{
  if(i&&!isNaN(i)&&i>0)
  {
    o.attr('class','green').html('+'+_num(i));
  }
  else
  {
    o.attr('class','red').html(_num(i));
  }
}
function _calc(p){
  var margin = (p*(data.conf.Margin/100)),
    sell = margin+p,
    amount = data.conf.Budget/p,
    profit = amount*margin,
    diff = (p*(data.conf.Cycle/100));
  return {'buy':p,'sell':sell,'margin':margin,'amount':amount,'profit':profit,'diff':diff};
}
function _sim(type,rate,amount)
{
  if(type=='sell')
  {
    var buy=rate/((data.conf.Margin+100)/100);
    var profit=(rate-buy)*amount;
    return {'do':'Buy: '+_num(_fs(buy)),'profit':_num(_fs(profit))+' '+currency}
  }
  else if(type=='buy')
  {
    var sell=((data.conf.Margin*rate)/100)+rate;
    var profit=(sell-rate)*amount;
    return {'do':'Sell: '+_num(_fs(sell)),'profit':_num(_fs(profit))+' '+currency}
  }
  return {'do':'','profit':''}
}
function _cycle(p){
  var mg = (data.conf.Cycle/100)*p,near=0;
  $.each(data.order,function(k,v){
    if(v.order_type=='sell')
    {
      if(Math.abs(v.rate-p)<mg)
      {
        near = v.rate;
      }
    }
  });
  return near;
}
function _num(x){
  if(!x)
  {
    return ''
  }
  var parts = x.toString().split(".");
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  return parts.join(".");
}
function _curr(t)
{
  $('#nav-curr .nav-link').removeClass('active');
  $('.nav-'+t).addClass('active');
  curr=t;
  if(tmg)
  {
    clearInterval(tmg);
    tmg='';
  }
}
var nav = {
  tmp:{},uparse:/^(((([^:\/#\?]+:)?(?:(\/\/)((?:(([^:@\/#\?]+)(?:\:([^:@\/#\?]+))?)@)?(([^:\/#\?\]\[]+|\[[^\/\]@#?]+\])(?:\:([0-9]+))?))?)?)?((\/?(?:[^\/\?#]+\/+)*)([^\?#]*)))?(\?[^#]+)?)(#.*)?/,
  load:function(){
    if(history&&history.pushState)
    {
      $(document).on('click','a',function(e){
        var t=$(this),href=t.attr('href'),target=t.attr('target'),pr=t.attr('data-pair');
        if((!e.ctrlKey)&&(!target)&&(pr)&&(href)&&(href.indexOf('javascript:')<0))nav.go(e,href)
      });
    }
  },
  go:function(e,href){
    var uri=nav.uparse.exec(window.location.href);
    var next=nav.uparse.exec(href);
    if(!next[13])next[13]='/';
    if(!uri[13])uri[13]='/';
    if((uri[1]==next[1])||((!next[11])&&(next[13]==uri[13])))
    {
      if(e)e.preventDefault();
      history.replaceState({url:next[13]}, null, href);
      network(next[13])
      return false;
    }
    else if((!next[11])||((next[11])&&(next[11]==uri[11])))
    {
      if(e)e.preventDefault();
      history.pushState({url:next[13]}, null, href);
      network(next[13])
      return false;
    }
    else if(!e)
    {
      console.log('not e');
      window.location.href=href;
    }
  },
  popstate:function(e)
  {
    console.log('popstate');
    console.log(e.originalEvent);
    nav.go(null,window.location.href);
  }
};
function order_type(){
  var type=$('.form-check-input[name=type]:checked').val();
  if(type=='buy'){
    $('#amount_type').html(data.pair.primary_currency);
  }else if(type=='sell'){
    $('#amount_type').html(data.pair.secondary_currency);
  }
}
function _getGraph(){
  if(!coin)return;
  jQuery.ajax({
    type: "GET",
    url: 'https://bxgo.io/ajax/graph?coin='+coin+'&callback=?',
    dataType: "jsonp",
    success: function(data){
      var d = new Date,v;
      var THB = data.rate.THB;
      var ar,av,ml=0,i=0,tm='';
      $.each(data.market,function(k,v){
        ar=[];
        av=v.Price||[];
        if(av.length>0){
          for(var j=0;j<av.length;j++){
            ar.push(_fs(av[j][3]*(curr=='thb'?THB:1)));
          }
          var price=_fs(v.Last[3]*(curr=='thb'?THB:1));
          ar.push(price);
          if(ar.length>ml)ml=ar.length;
          if(coin!=last_coin){
            last_coin=coin;
            window.LineConfig.data.datasets=[];
          }
          if(!window.LineConfig.data.datasets[i]){
            window.LineConfig.data.datasets.push({
                label: market[k]?market[k]:k+'.com',
                fill: false,
                backgroundColor: window.chartColors[i],
                borderColor: window.chartColors[i],
                data: ar,
            });
          }else{
            window.LineConfig.data.datasets[i].data=ar;
          }
          tm=v.Time;
          i++;
        }
      });
      var l=[];
      for(var i=0;i<ml-1;i++){
        l.push('');
      }
      l.push(tm);
      window.LineConfig.data.labels=l;
      window.LineChart.update();

    },
    error: function(XMLHttpRequest, textStatus, errorThrown){}
  });
}

$(window).on('popstate',function(e){nav.popstate(e);});
$(function(){
  tmr=setInterval(_getData,delay*1000)
  _getData();
  $('#order_modal').on('show.bs.modal',function(event){
    var btn = $(event.relatedTarget);
    $('#order_password').val(pass);
    $('#order_save').prop('checked',pass?true:false);
    $('#order_buysell').css('display','none');
    if(btn.data('action')=='cancel'){
      $('#order_title').html('ลบคำสั่งซื่อ/ขาย');
      $('#order_id').val(btn.data('id'));
    }
    else if(btn.data('action')=='config'){
      $('#order_title').html('Reload config.ini');
    }
    else if(btn.data('action')=='order'){
      $('#order_title').html('เพิ่มคำสั่งซื้อ/ขาย');
      $('#order_buysell').css('display','block');
      order_type();
    }
    $('#order_action').val(btn.data('action'));
  });
  nav.load();

  window.chartColors=['rgba(255,99,132,0.5)','rgba(255,159,64,0.5)','rgba(255,205,86,0.5)','rgba(75,192,192,0.5)','rgba(54,162,235,0.5)','rgba(153,102,255,0.5)'];
  window.LineConfig = {
      type: 'line',
      data: {
          labels:[],
          datasets: []
      },
      options: {
          responsive: true,
          maintainAspectRatio:false,
          elements:{point:{radius:0}},
          tooltips: {
              mode: 'index',
              intersect: false,
              callbacks: {
                  label: function(tooltipItem, data) {
                      return data.datasets[tooltipItem.datasetIndex].label+' : '+_fs(Number(tooltipItem.yLabel).toLocaleString('en'));
                  }
              }
          },
          hover: {
              mode: 'nearest',
              intersect: true
          },
          scales: {
            xAxes: [{
              display: true,
              gridLines:{display:true,drawBorder:true,drawOnChartArea:false,},
              ticks: {
                autoSkip: true,
                maxRotation: 0,
                minRotation: 0
              }
            }],
            yAxes: [{
              display: true,
              gridLines:{display:true,drawBorder:true,drawOnChartArea:false,},
              ticks: {
                callback: function (value) {
                  return Number(value).toLocaleString('en');
                },
              }
            }]
          }
      }
  };
  var ctx = document.getElementById("graph").getContext("2d");
  window.LineChart = new Chart(ctx, window.LineConfig);
});
