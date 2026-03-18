## Notater

### Ok heftig sjef, dette er spørsmål til studass siste timen:
* Kan man bare ha fuckton buffers på kanaler?
* Det var noe mer viktig, jeg må skru på huet.

## Hva mangler vi sjef
* Timeout på orders - check
* Obstruction solution - check
* Resend Assignments -
* Backup -
* Reconfiguratiom etter at en heis ble koblet fra (huske gamle caborders) 
* Process pair - 

Testing:
* Packet loss
* Resten av fat

## Tilbakemelding fra Sverre:
* Main challenge of primary-backup is configuring and reconfiguring the system

* How are we selecting primary and backup?\
 *Men ja dette må vi finne ut av*

 * Even the slaves send all data?\
 *No, slaves should only send new orders. And completed orders.*

 * Are you storing to disk?\
 *I guess no? I alle fall er local-storage veldig nedprioritert.*

 * You are making your own ack protocol?\
 *Ehh, idk, noen må jo si ifra at både primary og backupen har fått med seg alt.*


 * Make a /module/ configuring the system. Sending and receiving the
   heartbeats, and connecting/disconnecting the comm links, keeping
   track of primary and backup states, etc?\
   *I guess det er Peers og Reelection?*
 * Just use TCP for the rest, and send changes?\
 *Hmm*


## Hvordan designe ifølge sverre
**Motivating project discussion: Sverres Design Process**
1. Brainstorm for Use Cases: Span functionality space; do not aim for
"completeness".
2. Make design decisions: Divide into modules. (This is not a rational or
systematic process.)
3. Map the Use Cases from 1 on the design. This is both a cosistency
check and it yields the sub-use-cases for the modules.
4. Draw module interaction diagram. Who calls who?
5. Move responsibilities between modules (reorganize how the system is
divided into modules if necessary) so that the diagram in 4. gets fewer
arrows and that the module interfaces becomes perfect abstractions.
6. For each module: \
▶ Sum up use-cases from 3.\
▶ Either: Design the perfect module interface that satisfies the
use-cases - or recurse from 2.
7. Implement.



Brage som tester reelection på egen pc:

sudo ip route add 255.255.255.255 dev lo
sudo ip route del 255.255.255.255 dev lo