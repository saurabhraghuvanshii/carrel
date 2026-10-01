import java.util.*;

class RoomBooker {
    private final List<TreeMap<Integer, Integer>> booked = new ArrayList<>(); // per room: start -> end

    public RoomBooker(int rooms) {
        for (int i = 0; i < rooms; i++) {
            booked.add(new TreeMap<>());
        }
    }

    public int book(int start, int end) {
        for (int room = 0; room < booked.size(); room++) {
            TreeMap<Integer, Integer> b = booked.get(room);
            Map.Entry<Integer, Integer> before = b.lowerEntry(end);
            if (before == null || before.getValue() <= start) {
                b.put(start, end);
                return room;
            }
        }
        return -1;
    }

    public boolean cancel(int room, int start) {
        return booked.get(room).remove(start) != null;
    }
}
