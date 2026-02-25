## Notater

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