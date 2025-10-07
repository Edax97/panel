"""

"""

import time

import crcmod
import struct
import os
import time
import socket
import socketUtils
import crc16
import time
import random
import colored_logger
import logging
import colorama as cr
import sys
# python app_wialoncitos_coordinates.py 1415161998275621 67.172.246.251 5039
#crc16 = crcmod.mkCrcFun(0x13D65, 0xFFFF, True, 0xFFFF)
#  rpta = crc16(content)
logger = colored_logger.Logger("main", logging.DEBUG, cr.Fore.GREEN)
import route_generator


if "__main__" == __name__ :
    #get imei values
    if len(sys.argv) < 3:
        exit(1)

    for arg in sys.argv:
        print(arg)

    imei = sys.argv[1]
    ip = sys.argv[2]
    port = int(sys.argv[3])
    print(f"imei: {imei} ip: {ip} port: {port}")

    clientsocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    #clientsocket.connect(('127.0.0.1',5039))
    clientsocket.connect((ip, port))
    # Login   
    abc = crc16.crc16(f"2.0;{imei};NA;".encode())
    hex_value = f"{abc:x}"
    login_str = f"#L#2.0;{imei};NA;{hex_value}\r\n"
    print(login_str)
    clientsocket.send(login_str.encode())#(json.dumps(doc).encode()))
    rpta = socketUtils.recvData(clientsocket ,5, 3)
    if rpta:
        logger.info(rpta)
    else:
        logger.error("ERROR")
        exit(1)
    lat_dec = 0
    lon_dec = 0
    lat_int = random.randint(1191, 1222)# 1200  # top -11.912644, -76.987100 # bot -12.229742, -76.958260
    lon_int = random.randint(7500, 7800)# 7659
    while True:

        tiempo_segundos = int(time.time()) #- 18000  ## en linux esta gmt 00 en mi pc es gmt -5
        time_struct = time.localtime(tiempo_segundos)
        year= time.strftime('%y', time_struct)
        mes = time.strftime('%m', time_struct)
        dia = time.strftime('%d', time_struct)
        hora = time.strftime('%H', time_struct)
        minutos = time.strftime('%M', time_struct)
        sec = time.strftime('%S', time_struct)

        lat_dec += 20 
        lon_dec += 20
        if lat_dec == 1000:
            lat_dec = 0
            lat_int +=1
        if lon_dec == 1000:
            lat_dec = 0
            lat_int +=1
        sats = random.randint(5,20)
        course = random.randint(0, 360)
        #str = "{dia}{mes}{year};{hora}{minutos}{sec};{lat_int:0>4d}.{lat_dec:0>4d};S;0{lon_int:0>4d}.{lon_dec:0>4d};W;0;{course};265;{sats};NA;NA;NA;NA;NA;;".format(
        #    mes=mes, dia=dia , year=year, hora=hora, minutos= minutos, sec= sec, lat_int=lat_int, lat_dec = lat_dec, lon_int= lon_int, lon_dec = lon_dec, course = course, sats = sats
        #)
        #region    --- TESTING CUSTOM 
        #aceptado sin log ni lat ✅
        #str = "{dia}{mes}{year};{hora}{minutos}{sec};NA;NA;NA;NA;0;{course};265;{sats};NA;NA;NA;NA;NA;;".format(
        #    mes=mes, dia=dia , year=year, hora=hora, minutos= minutos, sec= sec, lat_int=lat_int, lat_dec = lat_dec, lon_int= lon_int, lon_dec = lon_dec, course = course, sats = sats
        #)
        # ni course ni altitud ni speed ni sats✅
        #str = "{dia}{mes}{year};{hora}{minutos}{sec};NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;".format(
        #    mes=mes, dia=dia , year=year, hora=hora, minutos= minutos, sec= sec, lat_int=lat_int, lat_dec = lat_dec, lon_int= lon_int, lon_dec = lon_dec, course = course, sats = sats
        #)
        # ✅
        str = "NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;fw:3:FLUJOv1.1.1,identificador:1:141516199827562,intensidad_wifi:1:-43,flujo_reverso_total:1:0,flujo_total:1:0,flujo_instantaneo:1:0,palabra_estado:1:0,horas_funcionamiento:1:0,nivel_bateria:1:0,frecuencia:1:15;".format(
            mes=mes, dia=dia , year=year, hora=hora, minutos= minutos, sec= sec, lat_int=lat_int, lat_dec = lat_dec, lon_int= lon_int, lon_dec = lon_dec, course = course, sats = sats
        )
        crc = crc16.crc16(str.encode())
        checksum16 = crc.to_bytes(2, "big").hex().upper()
        str = "#D#"+str+checksum16+ "\r\n"
        logger.info("Send> %s %s", imei, str)
        #endregion --- TESTING CUSTOM 


        #SD#NA;NA;NA;NA;NA;NA;NA;NA;NA;fw:3:FLUJOv1.1.1,identificador:1:141516199827562,intensidad_wifi:1:-43,flujo_reverso_total:1:0,flujo_total:1:0,flujo_instantaneo:1:0,palabra_estado:1:0,horas_funcionamiento:1:0,nivel_bateria:1:0,frecuencia:1:15;AEAB
        clientsocket.send(str.encode())
        
        
        rpta = socketUtils.recvData(clientsocket ,5, 3)
        logger.info(rpta)
        time.sleep(5)
