package fleet_test

import (
	"strconv"
	"strings"
	"testing"

	"example.com/vyshkat/internal/domain"
	"example.com/vyshkat/internal/fleet"
)

// TestEnergyLimitInspection фиксирует правило техосмотра по умолчанию (FR-3):
// расход больше лимита — отказ с причиной, ровно на лимите и без энергии — приём.
func TestEnergyLimitInspection(t *testing.T) {
	t.Parallel()

	sc := fleet.EnergyLimitInspection{MaxDailyKWh: fleet.DefaultMaxDailyKWh}
	under := must(domain.NewElectricScooter("VK-SC-0001", "Ninebot", 7, 0.4))(t)
	atLimit := must(domain.NewElectricBike("VK-EB-0001", "Eltreco", 6, 2.0))(t)
	over := must(domain.NewElectricBike("VK-EB-0002", "Old-Volt", 6, 2.01))(t)
	bike := must(domain.NewBicycle("VK-BK-0001", "Stels", 9))(t)

	tests := []struct {
		name         string
		vehicle      domain.Vehicle
		wantAccepted bool
		wantReason   string // фрагмент причины отказа; пусто для приёма
	}{
		{name: "расход меньше лимита", vehicle: under, wantAccepted: true},
		{name: "расход ровно на лимите", vehicle: atLimit, wantAccepted: true},
		{name: "расход больше лимита", vehicle: over, wantReason: "суточный расход 2.01 кВт·ч больше лимита 2.00"},
		{name: "велосипед без энергии", vehicle: bike, wantAccepted: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := sc.Inspect(tt.vehicle)

			if got.Accepted != tt.wantAccepted {
				t.Fatalf("Accepted = %v, ожидалось %v (причина: %q)", got.Accepted, tt.wantAccepted, got.Reason)
			}
			if got.Accepted != (got.Reason == "") {
				t.Errorf("Reason = %q не соответствует Accepted = %v", got.Reason, got.Accepted)
			}
			if !strings.Contains(got.Reason, tt.wantReason) {
				t.Errorf("Reason = %q, ожидался фрагмент %q", got.Reason, tt.wantReason)
			}
		})
	}
}

// TestMinSimplicityPolicy фиксирует правило «для новичков» на границе (FR-7).
func TestMinSimplicityPolicy(t *testing.T) {
	t.Parallel()

	policy := fleet.MinSimplicityPolicy{Min: fleet.DefaultMinSimplicity}

	tests := []struct {
		simplicity int
		want       bool
	}{
		{simplicity: domain.MinSimplicity, want: false},
		{simplicity: 5, want: false},
		{simplicity: 6, want: true},
		{simplicity: domain.MaxSimplicity, want: true},
	}

	for _, tt := range tests {
		t.Run("простота "+strconv.Itoa(tt.simplicity), func(t *testing.T) {
			t.Parallel()

			bike := must(domain.NewBicycle("VK-BK-0001", "Stels", tt.simplicity))(t)

			if got := policy.Allows(bike); got != tt.want {
				t.Errorf("Allows(простота %d) = %v, ожидалось %v", tt.simplicity, got, tt.want)
			}
		})
	}
}
