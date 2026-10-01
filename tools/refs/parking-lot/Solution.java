import java.util.*;

class ParkingLot {
    private final int small;
    private final PriorityQueue<Integer> freeSmall = new PriorityQueue<>();
    private final PriorityQueue<Integer> freeLarge = new PriorityQueue<>();
    private final Map<Integer, Integer> spotOf = new HashMap<>();

    public ParkingLot(int small, int large) {
        this.small = small;
        for (int s = 0; s < small; s++) {
            freeSmall.add(s);
        }
        for (int s = small; s < small + large; s++) {
            freeLarge.add(s);
        }
    }

    public int park(int car, int size) {
        if (spotOf.containsKey(car)) {
            return -1;
        }
        PriorityQueue<Integer> from = size == 1 && !freeSmall.isEmpty() ? freeSmall : freeLarge;
        if (from.isEmpty()) {
            return -1;
        }
        int spot = from.poll();
        spotOf.put(car, spot);
        return spot;
    }

    public int leave(int car) {
        Integer spot = spotOf.remove(car);
        if (spot == null) {
            return -1;
        }
        (spot < small ? freeSmall : freeLarge).add(spot);
        return spot;
    }

    public int freeSpots() {
        return freeSmall.size() + freeLarge.size();
    }
}
