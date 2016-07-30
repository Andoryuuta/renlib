#**WARING: (Work In Progress)** Extremely disorganized, lack of information, etc.

### <a name="preface"></a>Preface
C&C: Renegade development started approximately around 1998 ~ 1999, when 56k modems were in widespread use (According to wikipedia atleast). Because of the data restrictions, the protocol uses large amounts of bitpacking, delta compression, and floating point number scaling  

### <a name="notes"></a>Notes
The protocol *ONLY* uses UDP, any packets that need to be sent reliablely must use a proprietary TCP-esque method of acknowledging packets built on top of UPD.

In the tables below, "(BE)" represents Big Endian byte order and "(LE)" represents Little Endian byte order.

### <a name="constants"></a>Constants
```
#define MAX_PACKET_SIZE 0x224
```

# <a name="translayer"></a>Transmission layer
This section encompasses the stuff that is more associated with transmission of the data. Such as the packets moving straight from/to udp, packet grouping, and delta compresssion.


### <a name="udpmessage"></a>UDP Message
| 4 bytes (BE)  | [Up to [MAX_PACKET_SIZE](#constants) bytes)       |
|---------------|---------------------------------------|
| CRC-32 (IEEE) | [][Grouped packet](#groupedpacket) |

UDP messages are prefixed with a CRC32 checksum of the rest of the message data (the grouped packet data). 

### <a name="groupedpacket"></a>Grouped Packet
<table>
	<tr>
		<th>2 bytes (LE)</th>
		<th>header.Size bytes</th>
		<th>Variable amount of bytes</th>
	</tr>
  <tr>
    <td>Grouped Packet Header</td>
    <td>Packet Data</td>
    <td>[]Delta</td>
  </tr>
</table>

### <a name="groupedpacketheader"></a>Grouped Packet Header
<table>
	<tr>
		<th colspan="16">2 bytes (LE)</th>
	</tr>
	<tr>
		<th>15</th>
	  <th>14</th>
	  <th>13</th>
    <th>12</th>
    <th>11</th>
    <th>10</th>
    <th>9</th>
    <th>8</th>
    <th>7</th>
    <th>6</th>
    <th>5</th>
    <th>4</th>
    <th>3</th>
    <th>2</th>
    <th>1</th>
    <th>0</th>
  </tr>
  <tr>
    <td colspan="1" style="text-align: center">IsAnotherPacket</td>
    <td colspan="10" style="text-align: center">Size</td>
    <td colspan="5" style="text-align: center">Variations</td>
  </tr>
</table>

_IsAnotherPacket_ represents that there is another grouped packet after this one.<br />
_Size_ represents the size if the packet data in bytes<br />
_Variations_ represents how many variations of this packet there are. See [Delta Compression](#deltacompression)<br />

```Go
var header uint16
var IsAnotherPacket bool = ((header >> 15) & 0x01) == 1
var Size            uint = ((header >> 5) & 0x3ff)
var Variations      uint = ((header >> 0) & 0x1f)
```

### <a name="deltacompression"></a>Delta Compression
TODO: Write me!</br>
if _Variations_ > 1 then there are variations of the packet that use delta compression.

<table>
	<tr>
		<th>1 byte</th>
		<th>Variable amount of bytes</th>
		<th>Variable amount of bytes</th>
	</tr>
  <tr>
    <td>Header</td>
    <td>Bit masks/sets</td>
    <td>Delta diff data</td>
  </tr>
</table>



# <a name="applayer"></a>Application layer
The packet data after being ungrouped( and possibly deleta compressed) 
