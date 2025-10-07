import struct
import crcmod
import time
import select

import colored_logger
import logging
import colorama as cr

import selectors

logger = colored_logger.Logger("scoket Utils", logging.ERROR, cr.Fore.RED)


crc16 = crcmod.mkCrcFun(0x13D65, 0xFFFF, True, 0xFFFF)
OTA_SOF = 0xAA
OTA_EOF = 0xBB

OTA_PACKET_TYPE_CMD       = 0
OTA_PACKET_TYPE_DATA      = 1
OTA_PACKET_TYPE_HEADER    = 2
OTA_PACKET_TYPE_RESPONSE  = 3

OTA_CMD_START = 0
OTA_CMD_END   = 1
OTA_CMD_ABORT = 2

def recv_timeout(sock,bytes_to_read,timeout_seconds):
    sock.setblocking(0)
    ready = select.select([sock],[],[],timeout_seconds)
    if ready[0]:
        return sock.recv(bytes_to_read)

def recvData(sock, timeout, timeout_idle):
    inicio = time.time()
    sock.setblocking(False)
    sel = selectors.DefaultSelector()
    sel.register(sock, events=selectors.EVENT_READ, data = None)
    data = bytearray()
    while True:
        events = sel.select(timeout=timeout_idle)
        for key, mask in events:
            if mask & selectors.EVENT_READ:
                datarecv = key.fileobj.recv(1024)
                if datarecv:
                    data.extend(datarecv)
                    if data[-1:] == struct.pack("B", 0xBB):
                        return data 
                    if ( time.time() - inicio ) > timeout:
                        logger.error("periodo %d",  time.time() - inicio)
                        return data if len(data) > 0 else None
                else:
                    return data if len(data) > 0 else None
        if ( time.time() - inicio ) > timeout:
            logger.error("periodo %d",  time.time() - inicio)
            return data if len(data) > 0 else None
        
def recvData2(socket ,timeout, timeout_idle):
    inicio = time.time()
    data = bytearray()
    while True:
        try:
            packet = recv_timeout(socket, 4096, timeout_idle)
            if packet:
                data.extend(packet)
                inicio = time.time()
            if data[-1:] == struct.pack("B", 0xBB):
                #print("RECIBIDO...")
                return data
            if ( time.time() - inicio ) > timeout:
                print("periodo ",  time.time() - inicio)
                return data
            if len(packet) == 0:
                return None
        except Exception as err:        
            print("ERROR RECIBIENDO..."+str(err))
            return None

def recvDataCtrl(socket ,timeout, timeout_idle):
    inicio = time.time()
    data = bytearray()
    while True:
        if ( time.time() - inicio ) > timeout:
            #print("TIMEOUT RECIBIENDO...")
            if len(data)>2:
                return data
            return None
        try:
            packet = recv_timeout(socket, 4096, timeout_idle)
            if packet:
                data.extend(packet)
                inicio = time.time()
            if data[-2:] == b'\r\n':
                #print("RECIBIDO...")
                return data
                break;
        except Exception as err:        
            #print("ERROR RECIBIENDO..."+str(err))
            return None
    return None


def ota_cmd(ota_cmd):
    # sof: # ; packetType: 2b ; len: 2b ; Data: n bytes ; CRC: 4b ; EOF
    sof = struct.pack('B',OTA_SOF)#"#".encode('utf-8')
    packet_type = struct.pack('B', OTA_PACKET_TYPE_CMD)  # B unsigned char 1 byte
    data_len = struct.pack('H', 0) # H unsigned short 2 bytes
    
    #############################################################################
    #Data
    cmd = struct.pack('B',ota_cmd)
    data_len = struct.pack('H', 1)
    
    ##############################################################################
    #Content
    crc = struct.pack('I', crc16(bytes(cmd)))
    eof = struct.pack('B',OTA_EOF)
   
    return sof+ packet_type + data_len + cmd + crc + eof
    





def ota_header( filesize, filename):
    # sof: # ; packetType: 2b ; len: 2b ; HeaderData: 16b ; CRC: 4b ; EOF
    sof = struct.pack('B',OTA_SOF)
    packet_type = struct.pack('B', OTA_PACKET_TYPE_HEADER)  # B unsigned char 1 byte
    data_len = struct.pack('H', 0) # H unsigned short 2 bytes
    
    #############################################################################
    #Header data
    package_size = struct.pack('I', filesize)
    package_crc = struct.pack('I', 0)
    nada0 = struct.pack('I', 0)
    nada1 = struct.pack('I', 0)

    #if isinstance(filename, dict):
     #   filename = filename[1]
    if isinstance(filename, tuple):
        filename = filename[0]
    file_name = filename.encode()  #30 bytes
    if len(file_name) > 30:
        file_name = file_name[0:29] + b'\0'
    else: 
        file_name = file_name + bytes(30 - len(file_name))
    
    header_data =  package_size + package_crc + nada0 + nada1 + file_name
    data_len = struct.pack('H', len(header_data))
    
    ##############################################################################
    #Content
    crc = struct.pack('I', crc16(header_data))
    eof = struct.pack('B',OTA_EOF)
   
    return sof+ packet_type+ data_len + header_data + crc + eof
    
def ota_data(chunk):
    # sof: # ; packetType: 2b ; len: 2b ; Data: n bytes ; CRC: 4b ; EOF
    sof = struct.pack('B',OTA_SOF)
    packet_type = struct.pack('B', OTA_PACKET_TYPE_DATA)  # B unsigned char 1 byte
    data_len = struct.pack('H', 0) # H unsigned short 2 bytes
    
    #############################################################################
    #Data
    data = chunk
    data_len = struct.pack('H', len(data))
    
    ##############################################################################
    #Content
    crc = struct.pack('I', crc16(data))
    eof = struct.pack('B',OTA_EOF)
   
    return sof+ packet_type+ data_len + data + crc + eof
 
def ota_read_response(chunk):
    ##msg = struct.pack("BBHBIB", 0xAA,3,1,0,0,0xBB )   
    #msg = struct.unpack("BBHBIB", chunk )
    #print(f"rpta {chunk[4]}")
    if chunk[4] == 1:
        return False
    else:
        return True
        
        
#print("XXXXXX")
#makeLogin("Hola.ga", "65421" , "1.0", "L", "jajasaludos")


if __name__ == "__main__":
    pass