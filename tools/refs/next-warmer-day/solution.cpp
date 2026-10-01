#include <vector>
using namespace std;

vector<int> daysUntilWarmer(vector<int>& temps) {
    vector<int> wait(temps.size()), waiting;
    for (int i = 0; i < (int)temps.size(); i++) {
        while (!waiting.empty() && temps[waiting.back()] < temps[i]) {
            wait[waiting.back()] = i - waiting.back();
            waiting.pop_back();
        }
        waiting.push_back(i);
    }
    return wait;
}
