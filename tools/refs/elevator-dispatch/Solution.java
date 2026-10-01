import java.util.*;

class Elevator {
    private final TreeSet<Integer> stops = new TreeSet<>();
    private int at = 0, moved = 0;
    private boolean up = true;

    public Elevator(int floors) {
    }

    public void call(int floor) {
        if (floor != at) {
            stops.add(floor);
        }
    }

    public int next() {
        Integer stop = up ? stops.higher(at) : stops.lower(at);
        if (stop == null) {
            stop = up ? stops.lower(at) : stops.higher(at);
            if (stop == null) {
                return -1;
            }
            up = !up;
        }
        moved += Math.abs(stop - at);
        at = stop;
        stops.remove(stop);
        return stop;
    }

    public int travelled() {
        return moved;
    }
}
