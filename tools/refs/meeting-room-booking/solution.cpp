#include <map>
#include <vector>
using namespace std;

class RoomBooker {
    vector<map<int, int>> booked;  // per room: start -> end

public:
    RoomBooker(int rooms) : booked(rooms) {}

    int book(int start, int end) {
        for (int room = 0; room < (int)booked.size(); room++) {
            auto& b = booked[room];
            auto after = b.lower_bound(start);
            if (after != b.end() && after->first < end) continue;
            if (after != b.begin() && prev(after)->second > start) continue;
            b[start] = end;
            return room;
        }
        return -1;
    }

    bool cancel(int room, int start) { return booked[room].erase(start) > 0; }
};
