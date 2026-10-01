#include <functional>
#include <queue>
#include <unordered_map>
#include <vector>
using namespace std;

class ParkingLot {
    typedef priority_queue<int, vector<int>, greater<int>> MinHeap;
    int small;
    MinHeap freeSmall, freeLarge;
    unordered_map<int, int> spotOf;

public:
    ParkingLot(int small, int large) : small(small) {
        for (int s = 0; s < small; s++) freeSmall.push(s);
        for (int s = small; s < small + large; s++) freeLarge.push(s);
    }

    int park(int car, int size) {
        if (spotOf.count(car)) return -1;
        MinHeap& from = size == 1 && !freeSmall.empty() ? freeSmall : freeLarge;
        if (from.empty()) return -1;
        int spot = from.top();
        from.pop();
        spotOf[car] = spot;
        return spot;
    }

    int leave(int car) {
        auto it = spotOf.find(car);
        if (it == spotOf.end()) return -1;
        int spot = it->second;
        spotOf.erase(it);
        (spot < small ? freeSmall : freeLarge).push(spot);
        return spot;
    }

    int freeSpots() { return freeSmall.size() + freeLarge.size(); }
};
